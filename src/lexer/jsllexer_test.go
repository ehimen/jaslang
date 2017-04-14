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

func TestUnicode(t *testing.T) {
	doTestGetNext(t, "ϝЄ", []lexer.Lexeme{makeLexeme("ϝЄ", lexer.LIdentifier, 1)})
}

func TestSingleIdentifier(t *testing.T) {
	doTestGetNext(t, "foobar", []lexer.Lexeme{makeLexeme("foobar", lexer.LIdentifier, 1)})
}

func TestIdentifierAndString(t *testing.T) {
	doTestGetNext(
		t,
		"foo\"bar\"",
		[]lexer.Lexeme{
			makeLexeme("foo", lexer.LIdentifier, 1),
			makeLexeme("bar", lexer.LQuoted, 4),
		},
	)
}

func TestIdentWsString(t *testing.T) {
	doTestGetNext(
		t,
		"foo \"bar\"",
		[]lexer.Lexeme{
			makeLexeme("foo", lexer.LIdentifier, 1),
			makeLexeme(" ", lexer.LWhitespace, 4),
			makeLexeme("bar", lexer.LQuoted, 5),
		},
	)
}

func doTestGetNext(t *testing.T, in string, expected []lexer.Lexeme) {
	assertLexemes(
		t,
		lexer.NewJslLexer(getReader(in)),
		expected,
	)
}

func assertLexemes(t *testing.T, l lexer.Lexer, expected []lexer.Lexeme) {
	for _, s := range expected {
		current, err := l.GetNext()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
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
