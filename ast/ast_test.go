package ast

import (
	"github.com/lancelote/writing-an-interpreter-in-go/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	want := "let myVar = anotherVar;"
	got := program.String()

	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
