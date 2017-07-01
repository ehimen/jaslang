package testutil

import (
	"github.com/ehimen/jaslang/lex"
)

func MakeLexeme(value string, lexemeType lex.LexemeType, position int, line int) lex.Lexeme {
	return lex.Lexeme{
		Start: position,
		Type:  lexemeType,
		Value: value,
		Line:  line,
	}
}

// Generates a lexer implementation that simply returns
// the given lexemes in order, until it reaches the end
// where it return error EndOfInput.
func NewSimpleLexer(lexemes []lex.Lexeme) lex.Lexer {
	return &simpleLexer{0, lexemes}
}

type simpleLexer struct {
	current int
	lexemes []lex.Lexeme
}

func (l *simpleLexer) GetNext() (lex.Lexeme, error) {
	var lexeme lex.Lexeme

	if l.current >= len(l.lexemes) {
		return lexeme, lex.EndOfInput
	} else {
		defer func() { l.current++ }()

		return l.lexemes[l.current], nil
	}
}
