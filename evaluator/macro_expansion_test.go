package evaluator

import (
	"github.com/lancelote/writing-an-interpreter-in-go/ast"
	"github.com/lancelote/writing-an-interpreter-in-go/lexer"
	"github.com/lancelote/writing-an-interpreter-in-go/object"
	"github.com/lancelote/writing-an-interpreter-in-go/parser"
	"testing"
)

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1;
	let function = fn(x, y) { x + y };
	let mymacro = macro(x, y) { x + y; };
	`

	env := object.NewEnvironment()
	program := testParseProgram(input)

	DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("want 2 statements, got %d", len(program.Statements))
	}

	_, ok := env.Get("number")
	if ok {
		t.Fatalf("`number` should not be defined")
	}

	_, ok = env.Get("function")
	if ok {
		t.Fatalf("`function` should not be defined")
	}

	macroObj, ok := env.Get("mymacro")
	if !ok {
		t.Fatalf("macros is not in environment")
	}

	macro, ok := macroObj.(*object.Macro)
	if !ok {
		t.Fatalf("want macros, got %T", macroObj)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("want 2 macro parameters, got %d", len(macro.Parameters))
	}

	if macro.Parameters[0].String() != "x" {
		t.Fatalf("first macro parameter is not `x`, got %q", macro.Parameters[0].String())
	}

	if macro.Parameters[1].String() != "y" {
		t.Fatalf("second macro paramter is not `y`, got %q", macro.Parameters[1].String())
	}

	expectedBody := "(x + y)"

	if macro.Body.String() != expectedBody {
		t.Fatalf("want macro body %q, got %q", expectedBody, macro.Body.String())
	}
}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
