package lex

import "testing"

func TestVerbLexing(t *testing.T) {
	tests := map[string]Token{
		"POST":   Token{VERB, "POST"},
		"PATCH":  Token{VERB, "PATCH"},
		"DELETE": Token{VERB, "DELETE"},
		"GET":    Token{VERB, "GET"},
		"PUT":    Token{VERB, "PUT"},
	}

	for in, expected := range tests {
		lexer := NewLexer(in)
		actual, err := lexer.Next()
		if err != nil {
			t.Errorf(err)
			t.Fail()
		}
		if expected != actual {
			t.Fail()
		}
	}
}
