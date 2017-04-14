package lexer_test

import (
	"testing"
	"lexer"
	"strings"
)

func TestNewJslLexer(t *testing.T) {
	l := lexer.NewJslLexer(getReader(""))

	if _, ok := l.(lexer.Lexer); !ok {
		t.Error("Not a lexer")
	}
}

func TestGetNext(t *testing.T) {

	cases := []struct{
		in string
		out []lexer.Lexeme
	}{
		{
			"ϝЄ",
			[]lexer.Lexeme{makeLexeme("ϝЄ", lexer.LIdentifier, 1)},
		},
		{
			"foobar",
			[]lexer.Lexeme{
				makeLexeme("foobar", lexer.LIdentifier, 1),
			},
		},
		{
			"foo\"bar\"",
			[]lexer.Lexeme{
				makeLexeme("foo", lexer.LIdentifier, 1),
				makeLexeme("bar", lexer.LString, 4),
			},
		},
	}

	for _, testcase := range cases {
		assertLexemes(
			t,
			lexer.NewJslLexer(getReader(testcase.in)),
			testcase.out,
		)
	}
}

func assertLexemes(t *testing.T, l lexer.Lexer, expected []lexer.Lexeme) {
	for _, s := range expected {
		current, err := l.GetNext()

		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}

		if current != s {
			t.Errorf("Expected %s but got %s", s, current)
		}
	}

	if current, err := l.GetNext(); err != lexer.EndOfInput {
		t.Errorf("Expected end of input, but got %s", current)
	}
}

func getReader(s string) *strings.Reader {
	return strings.NewReader(s)
}

func makeLexeme(value string, lexemeType lexer.LexemeType, position int) lexer.Lexeme {
	return lexer.Lexeme{
		position,
		lexemeType,
		value,
	}
}
