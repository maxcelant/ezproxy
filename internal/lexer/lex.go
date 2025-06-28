package lex

type Lexer struct {
	// plaintext input source
	in string
	// Current position in the lexeme
	pos int
	// Start of lexeme
	start int
}

type StateFn func(*Lexer) StateFn

func NewLexer(in string) *Lexer {
	return &Lexer{
		in:    in,
		pos:   0,
		start: 0,
	}
}

func (l *Lexer) Next() (Token, error) {

}
