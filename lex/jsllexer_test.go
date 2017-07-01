package lex_test

import (
	"strings"
	"testing"

	"github.com/ehimen/jaslang/lex"
	"github.com/ehimen/jaslang/testutil"
)

func TestNewJslLexer(t *testing.T) {
	l := makeLexer("")

	if _, ok := l.(lex.Lexer); !ok {
		t.Error("Not a lexer")
	}
}

func TestUnicode(t *testing.T) {
	doTestGetNext(t, "ϝЄ", []lex.Lexeme{testutil.MakeLexeme("ϝЄ", lex.LIdentifier, 1, 1)})
}

func TestSingleIdentifier(t *testing.T) {
	doTestGetNext(t, "foobar", []lex.Lexeme{testutil.MakeLexeme("foobar", lex.LIdentifier, 1, 1)})
}

func TestIdentifierAndString(t *testing.T) {
	doTestGetNext(
		t,
		"foo\"bar\"",
		[]lex.Lexeme{
			testutil.MakeLexeme("foo", lex.LIdentifier, 1, 1),
			testutil.MakeLexeme("bar", lex.LQuoted, 4, 1),
		},
	)
}

func TestIdentWsString(t *testing.T) {
	doTestGetNext(
		t,
		"foo \"bar\"",
		[]lex.Lexeme{
			testutil.MakeLexeme("foo", lex.LIdentifier, 1, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 4, 1),
			testutil.MakeLexeme("bar", lex.LQuoted, 5, 1),
		},
	)
}

func TestOtherQuotedString(t *testing.T) {
	doTestGetNext(
		t,
		`'bar "foo"'`,
		[]lex.Lexeme{
			testutil.MakeLexeme(`bar "foo"`, lex.LQuoted, 1, 1),
		},
	)
}

func TestOtherOtherQuotedString(t *testing.T) {
	doTestGetNext(
		t,
		`"bar 'foo'"`,
		[]lex.Lexeme{
			testutil.MakeLexeme("bar 'foo'", lex.LQuoted, 1, 1),
		},
	)
}

func TestEscapedQuote(t *testing.T) {
	doTestGetNext(
		t,
		`"\""`,
		[]lex.Lexeme{
			testutil.MakeLexeme(`"`, lex.LQuoted, 1, 1),
		},
	)
}

func TestEscapedBackslash(t *testing.T) {
	doTestGetNext(
		t,
		`"\\"`,
		[]lex.Lexeme{
			testutil.MakeLexeme(`\`, lex.LQuoted, 1, 1),
		},
	)
}

func TestEscapedOtherQuote(t *testing.T) {
	doTestGetNext(
		t,
		`"\'"`,
		[]lex.Lexeme{
			testutil.MakeLexeme(`\'`, lex.LQuoted, 1, 1),
		},
	)
}

func TestBackslashesAndEscapedQuotes(t *testing.T) {
	doTestGetNext(
		t,
		`"\\\""`,
		[]lex.Lexeme{
			testutil.MakeLexeme(`\"`, lex.LQuoted, 1, 1),
		},
	)
}

func TestUnterminatedString(t *testing.T) {
	lexer := makeLexer(`"foo`)
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
			testutil.MakeLexeme("{", lex.LBraceOpen, 1, 1),
			testutil.MakeLexeme("}", lex.LBraceClose, 2, 1),
			testutil.MakeLexeme("(", lex.LParenOpen, 3, 1),
			testutil.MakeLexeme(")", lex.LParenClose, 4, 1),
			testutil.MakeLexeme(";", lex.LSemiColon, 5, 1),
		},
	)
}

func TestNumberSimple(t *testing.T) {
	doTestGetNext(
		t,
		"1",
		[]lex.Lexeme{
			testutil.MakeLexeme("1", lex.LNumber, 1, 1),
		},
	)
}

func TestNumbersSigned(t *testing.T) {
	doTestGetNext(
		t,
		"-1+2",
		[]lex.Lexeme{
			testutil.MakeLexeme("-1", lex.LNumber, 1, 1),
			testutil.MakeLexeme("+2", lex.LNumber, 3, 1),
		},
	)
}

func TestDecimal(t *testing.T) {
	doTestGetNext(
		t,
		"1.34",
		[]lex.Lexeme{
			testutil.MakeLexeme("1.34", lex.LNumber, 1, 1),
		},
	)
}

func TestNegativeDecimal(t *testing.T) {
	doTestGetNext(
		t,
		"-1.34",
		[]lex.Lexeme{
			testutil.MakeLexeme("-1.34", lex.LNumber, 1, 1),
		},
	)
}

func TestSingleSign(t *testing.T) {
	doTestGetNext(
		t,
		"+",
		[]lex.Lexeme{
			testutil.MakeLexeme("+", lex.LOperator, 1, 1),
		},
	)
}

func TestMultipleDecimalPoints(t *testing.T) {
	doTestGetNext(
		t,
		"13.14.24",
		[]lex.Lexeme{
			testutil.MakeLexeme("13.14", lex.LNumber, 1, 1),
			testutil.MakeLexeme(".", lex.LOperator, 6, 1),
			testutil.MakeLexeme("24", lex.LNumber, 7, 1),
		},
	)
}

func TestSingleWidthOperators(t *testing.T) {
	doTestGetNext(
		t,
		"+ -",
		[]lex.Lexeme{
			testutil.MakeLexeme("+", lex.LOperator, 1, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 2, 1),
			testutil.MakeLexeme("-", lex.LOperator, 3, 1),
		},
	)
}

func TestOperatorOverIdentifier(t *testing.T) {
	doTestGetNext(
		t,
		"+foo",
		[]lex.Lexeme{
			testutil.MakeLexeme("+", lex.LOperator, 1, 1),
			testutil.MakeLexeme("foo", lex.LIdentifier, 2, 1),
		},
	)
}

func TestFunctionCall(t *testing.T) {
	doTestGetNext(
		t,
		"foo(bar)",
		[]lex.Lexeme{
			testutil.MakeLexeme("foo", lex.LIdentifier, 1, 1),
			testutil.MakeLexeme("(", lex.LParenOpen, 4, 1),
			testutil.MakeLexeme("bar", lex.LIdentifier, 5, 1),
			testutil.MakeLexeme(")", lex.LParenClose, 8, 1),
		},
	)
}

func TestIf(t *testing.T) {
	doTestGetNext(
		t,
		"if(){}elseif{}else{}",
		[]lex.Lexeme{
			testutil.MakeLexeme("if", lex.LIf, 1, 1),
			testutil.MakeLexeme("(", lex.LParenOpen, 3, 1),
			testutil.MakeLexeme(")", lex.LParenClose, 4, 1),
			testutil.MakeLexeme("{", lex.LBraceOpen, 5, 1),
			testutil.MakeLexeme("}", lex.LBraceClose, 6, 1),
			testutil.MakeLexeme("elseif", lex.LElseIf, 7, 1),
			testutil.MakeLexeme("{", lex.LBraceOpen, 13, 1),
			testutil.MakeLexeme("}", lex.LBraceClose, 14, 1),
			testutil.MakeLexeme("else", lex.LElse, 15, 1),
			testutil.MakeLexeme("{", lex.LBraceOpen, 19, 1),
			testutil.MakeLexeme("}", lex.LBraceClose, 20, 1),
		},
	)
}

func TestLet(t *testing.T) {
	doTestGetNext(
		t,
		"let a = 1",
		[]lex.Lexeme{
			testutil.MakeLexeme("let", lex.LLet, 1, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 4, 1),
			testutil.MakeLexeme("a", lex.LIdentifier, 5, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 6, 1),
			testutil.MakeLexeme("=", lex.LEquals, 7, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 8, 1),
			testutil.MakeLexeme("1", lex.LNumber, 9, 1),
		},
	)
}

func TestWhile(t *testing.T) {
	doTestGetNext(
		t,
		"while(){}",
		[]lex.Lexeme{
			testutil.MakeLexeme("while", lex.LWhile, 1, 1),
			testutil.MakeLexeme("(", lex.LParenOpen, 6, 1),
			testutil.MakeLexeme(")", lex.LParenClose, 7, 1),
			testutil.MakeLexeme("{", lex.LBraceOpen, 8, 1),
			testutil.MakeLexeme("}", lex.LBraceClose, 9, 1),
		},
	)
}

func TestBoolValues(t *testing.T) {
	doTestGetNext(
		t,
		"true false TRUE FALSE",
		[]lex.Lexeme{
			testutil.MakeLexeme("true", lex.LBoolTrue, 1, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 5, 1),
			testutil.MakeLexeme("false", lex.LBoolFalse, 6, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 11, 1),
			testutil.MakeLexeme("TRUE", lex.LBoolTrue, 12, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 16, 1),
			testutil.MakeLexeme("FALSE", lex.LBoolFalse, 17, 1),
		},
	)
}

func TestMultipleWhitespace(t *testing.T) {
	doTestGetNext(
		t,
		" \t\n",
		[]lex.Lexeme{
			testutil.MakeLexeme(" \t\n", lex.LWhitespace, 1, 1),
		},
	)
}

func TestIdentifierWithSpecialPrefix(t *testing.T) {
	doTestGetNext(
		t,
		"iffoo",
		[]lex.Lexeme{
			testutil.MakeLexeme("iffoo", lex.LIdentifier, 1, 1),
		},
	)
}

func TestOperatorBetweenIdentifiers(t *testing.T) {
	doTestGetNext(
		t,
		"foo+bar",
		[]lex.Lexeme{
			testutil.MakeLexeme("foo", lex.LIdentifier, 1, 1),
			testutil.MakeLexeme("+", lex.LOperator, 4, 1),
			testutil.MakeLexeme("bar", lex.LIdentifier, 5, 1),
		},
	)
}

func TestNotSingleOperatorBetweenIdentifiers(t *testing.T) {
	doTestGetNext(
		t,
		"foo++bar",
		[]lex.Lexeme{
			testutil.MakeLexeme("foo", lex.LIdentifier, 1, 1),
			testutil.MakeLexeme("++", lex.LOperator, 4, 1),
			testutil.MakeLexeme("bar", lex.LIdentifier, 6, 1),
		},
	)
}

func TestFunctionDeclaration(t *testing.T) {
	doTestGetNext(
		t,
		"() => {}",
		[]lex.Lexeme{
			testutil.MakeLexeme("(", lex.LParenOpen, 1, 1),
			testutil.MakeLexeme(")", lex.LParenClose, 2, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 3, 1),
			testutil.MakeLexeme("=>", lex.LOperator, 4, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 6, 1),
			testutil.MakeLexeme("{", lex.LBraceOpen, 7, 1),
			testutil.MakeLexeme("}", lex.LBraceClose, 8, 1),
		},
	)
}

func TestGte(t *testing.T) {
	doTestGetNext(
		t,
		"3>=4",
		[]lex.Lexeme{
			testutil.MakeLexeme("3", lex.LNumber, 1, 1),
			testutil.MakeLexeme(">=", lex.LOperator, 2, 1),
			testutil.MakeLexeme("4", lex.LNumber, 4, 1),
		},
	)
}

func TestAssignment(t *testing.T) {
	code :=
		`let n number = 4;
let total number = 0;`

	doTestGetNext(
		t,
		code,
		[]lex.Lexeme{
			testutil.MakeLexeme("let", lex.LLet, 1, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 4, 1),
			testutil.MakeLexeme("n", lex.LIdentifier, 5, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 6, 1),
			testutil.MakeLexeme("number", lex.LIdentifier, 7, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 13, 1),
			testutil.MakeLexeme("=", lex.LEquals, 14, 1),
			testutil.MakeLexeme(" ", lex.LWhitespace, 15, 1),
			testutil.MakeLexeme("4", lex.LNumber, 16, 1),
			testutil.MakeLexeme(";", lex.LSemiColon, 17, 1),
			testutil.MakeLexeme("\n", lex.LWhitespace, 18, 1),
			testutil.MakeLexeme("let", lex.LLet, 1, 2),
			testutil.MakeLexeme(" ", lex.LWhitespace, 4, 2),
			testutil.MakeLexeme("total", lex.LIdentifier, 5, 2),
			testutil.MakeLexeme(" ", lex.LWhitespace, 10, 2),
			testutil.MakeLexeme("number", lex.LIdentifier, 11, 2),
			testutil.MakeLexeme(" ", lex.LWhitespace, 17, 2),
			testutil.MakeLexeme("=", lex.LEquals, 18, 2),
			testutil.MakeLexeme(" ", lex.LWhitespace, 19, 2),
			testutil.MakeLexeme("0", lex.LNumber, 20, 2),
			testutil.MakeLexeme(";", lex.LSemiColon, 21, 2),
		},
	)
}

func doTestGetNext(t *testing.T, in string, expected []lex.Lexeme) {
	assertLexemes(
		t,
		makeLexer(in),
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

func makeLexer(input string) lex.Lexer {
	return lex.NewJslLexer(getReader(input))
}
