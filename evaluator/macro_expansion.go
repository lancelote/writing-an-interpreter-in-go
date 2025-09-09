package evaluator

import (
	"github.com/lancelote/writing-an-interpreter-in-go/ast"
	"github.com/lancelote/writing-an-interpreter-in-go/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, stmt := range program.Statements {
		if isMacroDefinition(stmt) {
			addMacro(stmt, env)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i-- {
		definitionIdx := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIdx],
			program.Statements[definitionIdx+1:]...,
		)
	}
}

func isMacroDefinition(node ast.Statement) bool {
	letStmt, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStmt.Value.(*ast.MacroLiteral)
	if !ok {
		return false
	}

	return true
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	letStmt, _ := stmt.(*ast.LetStatement)
	macroLiteral, _ := letStmt.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Env:        env,
		Body:       macroLiteral.Body,
	}

	env.Set(letStmt.Name.Value, macro)
}
