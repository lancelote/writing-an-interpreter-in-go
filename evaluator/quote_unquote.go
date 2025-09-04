package evaluator

import (
	"github.com/lancelote/writing-an-interpreter-in-go/ast"
	"github.com/lancelote/writing-an-interpreter-in-go/object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
