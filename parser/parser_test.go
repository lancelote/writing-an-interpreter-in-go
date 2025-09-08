package parser

import (
	"fmt"
	"github.com/lancelote/writing-an-interpreter-in-go/ast"
	"github.com/lancelote/writing-an-interpreter-in-go/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParseErrors(t, p)

		assertStatementCount(t, program.Statements, 1)

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input               string
		expectedReturnValue any
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		assertStatementCount(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("want return statement, got %T", stmt)
		}

		if !testLiteralExpression(t, stmt.ReturnValue, tt.expectedReturnValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

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

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

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

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("want if-expression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	assertStatementCount(t, exp.Consequence.Statements, 1)

	consequence := assertExpressionStatement(t, exp.Consequence.Statements[0])

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("alternative is not nil, got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("want if-expression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	assertStatementCount(t, exp.Consequence.Statements, 1)

	consequence := assertExpressionStatement(t, exp.Consequence.Statements[0])

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	assertStatementCount(t, exp.Alternative.Statements, 1)

	alternative := assertExpressionStatement(t, exp.Alternative.Statements[0])

	if !testIdentifier(t, alternative.Expression, "y") {
		return
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

		stmt := assertExpressionStatement(t, program.Statements[0])

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

		stmt := assertExpressionStatement(t, program.Statements[0])

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
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
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

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("want function literal, got %T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("unexpected number of parameters, want 2, got %d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	assertStatementCount(t, function.Body.Statements, 1)

	bodyStmt := assertExpressionStatement(t, function.Body.Statements[0])

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		stmt := assertExpressionStatement(t, program.Statements[0])
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("want %d parameters, got %d", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("want call expression, got %T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("want 3 arguments, got %d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionArgumentParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{"add();", "add", []string{}},
		{"add(1);", "add", []string{"1"}},
		{"add(1, 2 * 3, 4 + 5);", "add", []string{"1", "(2 * 3)", "(4 + 5)"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		assertStatementCount(t, program.Statements, 1)
		stmt := assertExpressionStatement(t, program.Statements[0])

		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("want call expression, got %T", stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(tt.expectedArgs) != len(exp.Arguments) {
			t.Fatalf("want %d arguments, got %d", len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if arg != exp.Arguments[i].String() {
				t.Errorf("want argument %s, got %s", arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal, got %T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("expected %q string literal, got %q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("want array literal, got %T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("want 3 array elements, got %d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingEmptyArrayLiteral(t *testing.T) {
	input := "[]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("want array literal, got %T", stmt.Expression)
	}

	if len(array.Elements) != 0 {
		t.Fatalf("want empty array, got length %d", len(array.Elements))
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("want index expression, got %T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("want hash literal, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("want 3 pairs, got %d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("want string key, got %T", key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := `{}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("want hash literal, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("want 0 pairs, got %d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWitExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("want hash literal, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("want 3 pairs, got %d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("want string key, got %T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("no test function for key %q found", literal.String())
		}

		testFunc(value)
	}
}

func TestMacroLiteralParsing(t *testing.T) {
	input := `macro(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	assertStatementCount(t, program.Statements, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	macro, ok := stmt.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("want macro literal, got %T", stmt)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("want 2 macro paramters, got %d", len(macro.Parameters))
	}

	testLiteralExpression(t, macro.Parameters[0], "x")
	testLiteralExpression(t, macro.Parameters[1], "y")

	assertStatementCount(t, macro.Body.Statements, 1)

	bodyStmt := assertExpressionStatement(t, macro.Body.Statements[0])

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
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
		t.Errorf("expected identifier value %s, got %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("expected identifier token literal %s, got %s", value, ident.TokenLiteral())
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

func assertStatementCount(t *testing.T, stmts []ast.Statement, expected int) {
	t.Helper()

	if len(stmts) != expected {
		t.Fatalf("got %d statements, want %d", len(stmts), expected)
	}
}

func assertExpressionStatement(t *testing.T, obj any) *ast.ExpressionStatement {
	t.Helper()

	v, ok := obj.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("want expression statement, got=%T", obj)
	}
	return v
}
