package lex_test

import (
	"lex"
	"strings"
	"testing"
)

func TestNewJslLexer(t *testing.T) {
	l := makeOperatorlessLexer("")

	if _, ok := l.(lex.Lexer); !ok {
		t.Error("Not a lexer")
	}
}

func TestUnicode(t *testing.T) {
	doTestGetNext(t, "ϝЄ", []lex.Lexeme{makeLexeme("ϝЄ", lex.LIdentifier, 1)})
}

func TestSingleIdentifier(t *testing.T) {
	doTestGetNext(t, "foobar", []lex.Lexeme{makeLexeme("foobar", lex.LIdentifier, 1)})
}

func TestIdentifierAndString(t *testing.T) {
	doTestGetNext(
		t,
		"foo\"bar\"",
		[]lex.Lexeme{
			makeLexeme("foo", lex.LIdentifier, 1),
			makeLexeme("bar", lex.LQuoted, 4),
		},
	)
}

func TestIdentWsString(t *testing.T) {
	doTestGetNext(
		t,
		"foo \"bar\"",
		[]lex.Lexeme{
			makeLexeme("foo", lex.LIdentifier, 1),
			makeLexeme(" ", lex.LWhitespace, 4),
			makeLexeme("bar", lex.LQuoted, 5),
		},
	)
}

func TestOtherQuotedString(t *testing.T) {
	doTestGetNext(
		t,
		`'bar "foo"'`,
		[]lex.Lexeme{
			makeLexeme(`bar "foo"`, lex.LQuoted, 1),
		},
	)
}

func TestOtherOtherQuotedString(t *testing.T) {
	doTestGetNext(
		t,
		`"bar 'foo'"`,
		[]lex.Lexeme{
			makeLexeme("bar 'foo'", lex.LQuoted, 1),
		},
	)
}

func TestEscapedQuote(t *testing.T) {
	doTestGetNext(
		t,
		`"\""`,
		[]lex.Lexeme{
			makeLexeme(`"`, lex.LQuoted, 1),
		},
	)
}

func TestEscapedBackslash(t *testing.T) {
	doTestGetNext(
		t,
		`"\\"`,
		[]lex.Lexeme{
			makeLexeme(`\`, lex.LQuoted, 1),
		},
	)
}

func TestEscapedOtherQuote(t *testing.T) {
	doTestGetNext(
		t,
		`"\'"`,
		[]lex.Lexeme{
			makeLexeme(`\'`, lex.LQuoted, 1),
		},
	)
}

func TestBackslashesAndEscapedQuotes(t *testing.T) {
	doTestGetNext(
		t,
		`"\\\""`,
		[]lex.Lexeme{
			makeLexeme(`\"`, lex.LQuoted, 1),
		},
	)
}

func TestUnterminatedString(t *testing.T) {
	lexer := makeOperatorlessLexer(`"foo`)
	_, err := lexer.GetNext()

	if err != lex.UnterminatedString {
		t.Errorf("Expected unterminated string, but got %v", err)
	}
}

func TestCharacterSymbols(t *testing.T) {
	doTestGetNext(
		t,
		"{}();",
		[]lex.Lexeme{
			makeLexeme("{", lex.LBraceOpen, 1),
			makeLexeme("}", lex.LBraceClose, 2),
			makeLexeme("(", lex.LParenOpen, 3),
			makeLexeme(")", lex.LParenClose, 4),
			makeLexeme(";", lex.LSemiColon, 5),
		},
	)
}

func TestNumberSimple(t *testing.T) {
	doTestGetNext(
		t,
		"1",
		[]lex.Lexeme{
			makeLexeme("1", lex.LNumber, 1),
		},
	)
}

func TestNumbersSigned(t *testing.T) {
	doTestGetNext(
		t,
		"-1+2",
		[]lex.Lexeme{
			makeLexeme("-1", lex.LNumber, 1),
			makeLexeme("+2", lex.LNumber, 3),
		},
	)
}

func TestDecimal(t *testing.T) {
	doTestGetNext(
		t,
		"1.34",
		[]lex.Lexeme{
			makeLexeme("1.34", lex.LNumber, 1),
		},
	)
}

func TestNegativeDecimal(t *testing.T) {
	doTestGetNext(
		t,
		"-1.34",
		[]lex.Lexeme{
			makeLexeme("-1.34", lex.LNumber, 1),
		},
	)
}

func TestSingleSign(t *testing.T) {
	doTestGetNext(
		t,
		"+",
		[]lex.Lexeme{
			makeLexeme("+", lex.LIdentifier, 1),
		},
	)
}

func TestMultipleDecimalPoints(t *testing.T) {
	doTestGetNext(
		t,
		"13.14.24",
		[]lex.Lexeme{
			makeLexeme("13.14", lex.LNumber, 1),
			makeLexeme(".24", lex.LIdentifier, 6),
		},
	)
}

func TestSingleWidthOperators(t *testing.T) {
	doTestGetNextWithOperators(
		t,
		"+ -",
		[]string{"+", "-"},
		[]lex.Lexeme{
			makeLexeme("+", lex.LOperator, 1),
			makeLexeme(" ", lex.LWhitespace, 2),
			makeLexeme("-", lex.LOperator, 3),
		},
	)
}

func TestOperatorOverIdentifier(t *testing.T) {
	doTestGetNextWithOperators(
		t,
		"+foo",
		[]string{"+"},
		[]lex.Lexeme{
			makeLexeme("+", lex.LOperator, 1),
			makeLexeme("foo", lex.LIdentifier, 2),
		},
	)
}

func TestFunctionCall(t *testing.T) {
	doTestGetNextWithOperators(
		t,
		"foo(bar)",
		[]string{"+", "-"},
		[]lex.Lexeme{
			makeLexeme("foo", lex.LIdentifier, 1),
			makeLexeme("(", lex.LParenOpen, 4),
			makeLexeme("bar", lex.LIdentifier, 5),
			makeLexeme(")", lex.LParenClose, 8),
		},
	)
}

func doTestGetNext(t *testing.T, in string, expected []lex.Lexeme) {
	assertLexemes(
		t,
		makeOperatorlessLexer(in),
		expected,
	)
}

func doTestGetNextWithOperators(t *testing.T, in string, operators []string, expected []lex.Lexeme) {
	assertLexemes(
		t,
		makeLexer(in, operators),
		expected,
	)
}

func assertLexemes(t *testing.T, l lex.Lexer, expected []lex.Lexeme) {
	for _, s := range expected {
		current, err := l.GetNext()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if current != s {
			t.Errorf("Expected %s but got %s", s, current)
		}
	}

	if current, err := l.GetNext(); err != lex.EndOfInput {
		t.Errorf("Expected end of input, but got %s", current)
	}
}

func getReader(s string) *strings.Reader {
	return strings.NewReader(s)
}

func makeLexer(input string, operators []string) lex.Lexer {
	return lex.NewJslLexer(getReader(input), operators)
}

func makeOperatorlessLexer(input string) lex.Lexer {
	return lex.NewJslLexer(getReader(input), []string{})
}

func makeLexeme(value string, lexemeType lex.LexemeType, position int) lex.Lexeme {
	return lex.Lexeme{
		position,
		lexemeType,
		value,
	}
}
