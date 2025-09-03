package evaluator

import (
	"github.com/lancelote/writing-an-interpreter-in-go/object"
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(5)`,
			`5`,
		},
		{
			`quote(5 + 8)`,
			`(5 + 8)`,
		},
		{
			`quote(foobar + barfoo)`,
			`(foobar + barfoo)`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected quote object, got %T", evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote object is nil")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("want %q, got %q", tt.expected, quote.Node.String())
		}
	}
}
