package parser

import (
	"fmt"
	"github.com/lancelote/writing-an-interpreter-in-go/ast"
	"github.com/lancelote/writing-an-interpreter-in-go/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	assertStatementCount(t, program, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	assertStatementCount(t, program, 3)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not *ast.ReturnStatement, got %T", stmt)
			continue
		}
		literal := returnStmt.TokenLiteral()
		if literal != "return" {
			t.Errorf("returnStmt literal is not `return`, got %q", literal)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assertStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("want expression statement, got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("want identifier, got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("expected identifier value `foobar`, got=`%s`", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("expected identifier token literal `foobar`, got=`%s`", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assertStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("want expression statement, got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("want integer literal, got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("expected integer literal value `5`, got=`%d`", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("expected identifier token literal `5`, got=`%s`", literal.TokenLiteral())
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assertStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("want expression statement, got=%T", program.Statements[0])
	}

	boolean, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("want boolean, got=%T", stmt.Expression)
	}

	if boolean.Value != true {
		t.Errorf("expected `true`, got `%t`", boolean.Value)
	}

	if boolean.TokenLiteral() != "true" {
		t.Errorf("expected boolean token literal `true`, got=`%s`", boolean.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("not enough statements, got=%d, want 1", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("want expression statement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("want prefix expression, got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not `%s`, got=`%s`", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if count := len(program.Statements); count != 1 {
			t.Fatalf("not enough statements, got=%d, want 1", count)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("want expression statement, got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}

	t.Errorf("type of expression is not handled, got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expected integer literal, got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("expected integer %d, got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("token literal expected %d, got=%s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("want identifier, got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("expected identifier value %s, got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("expected identifier token literal %s, got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("expected boolean, got=%T", exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean value is not %t, got %t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("expected boolean token literal %t, got %s", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("want infix expression, got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("%s operator expected, got %s", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func checkParseErrors(t *testing.T, p *Parser) {
	t.Helper()

	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() is not 'let', got '%q'", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement, got '%T'", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not '%s', got '%s'", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s', got '%s'", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func assertStatementCount(t *testing.T, program *ast.Program, expected int) {
	t.Helper()

	if len(program.Statements) != expected {
		t.Fatalf("got %d statements, want %d", len(program.Statements), expected)
	}
}
