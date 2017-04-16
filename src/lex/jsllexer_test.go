package lex_test

import (
	"lex"
	"strings"
	"testing"
)

func TestNewJslLexer(t *testing.T) {
	l := makeLexer("")

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
			makeLexeme("+", lex.LOperator, 1),
		},
	)
}

func TestMultipleDecimalPoints(t *testing.T) {
	doTestGetNext(
		t,
		"13.14.24",
		[]lex.Lexeme{
			makeLexeme("13.14", lex.LNumber, 1),
			makeLexeme(".", lex.LOperator, 6),
			makeLexeme("24", lex.LNumber, 7),
		},
	)
}

func TestSingleWidthOperators(t *testing.T) {
	doTestGetNext(
		t,
		"+ -",
		[]lex.Lexeme{
			makeLexeme("+", lex.LOperator, 1),
			makeLexeme(" ", lex.LWhitespace, 2),
			makeLexeme("-", lex.LOperator, 3),
		},
	)
}

func TestOperatorOverIdentifier(t *testing.T) {
	doTestGetNext(
		t,
		"+foo",
		[]lex.Lexeme{
			makeLexeme("+", lex.LOperator, 1),
			makeLexeme("foo", lex.LIdentifier, 2),
		},
	)
}

func TestFunctionCall(t *testing.T) {
	doTestGetNext(
		t,
		"foo(bar)",
		[]lex.Lexeme{
			makeLexeme("foo", lex.LIdentifier, 1),
			makeLexeme("(", lex.LParenOpen, 4),
			makeLexeme("bar", lex.LIdentifier, 5),
			makeLexeme(")", lex.LParenClose, 8),
		},
	)
}

func TestIf(t *testing.T) {
	doTestGetNext(
		t,
		"if(){}elseif{}else{}",
		[]lex.Lexeme{
			makeLexeme("if", lex.LIf, 1),
			makeLexeme("(", lex.LParenOpen, 3),
			makeLexeme(")", lex.LParenClose, 4),
			makeLexeme("{", lex.LBraceOpen, 5),
			makeLexeme("}", lex.LBraceClose, 6),
			makeLexeme("elseif", lex.LElseIf, 7),
			makeLexeme("{", lex.LBraceOpen, 13),
			makeLexeme("}", lex.LBraceClose, 14),
			makeLexeme("else", lex.LElse, 15),
			makeLexeme("{", lex.LBraceOpen, 19),
			makeLexeme("}", lex.LBraceClose, 20),
		},
	)
}

func TestLet(t *testing.T) {
	doTestGetNext(
		t,
		"let a = 1",
		[]lex.Lexeme{
			makeLexeme("let", lex.LLet, 1),
			makeLexeme(" ", lex.LWhitespace, 4),
			makeLexeme("a", lex.LIdentifier, 5),
			makeLexeme(" ", lex.LWhitespace, 6),
			makeLexeme("=", lex.LOperator, 7),
			makeLexeme(" ", lex.LWhitespace, 8),
			makeLexeme("1", lex.LNumber, 9),
		},
	)
}

func TestWhile(t *testing.T) {
	doTestGetNext(
		t,
		"while(){}",
		[]lex.Lexeme{
			makeLexeme("while", lex.LWhile, 1),
			makeLexeme("(", lex.LParenOpen, 6),
			makeLexeme(")", lex.LParenClose, 7),
			makeLexeme("{", lex.LBraceOpen, 8),
			makeLexeme("}", lex.LBraceClose, 9),
		},
	)
}

func TestBoolValues(t *testing.T) {
	doTestGetNext(
		t,
		"true false TRUE FALSE",
		[]lex.Lexeme{
			makeLexeme("true", lex.LBoolTrue, 1),
			makeLexeme(" ", lex.LWhitespace, 5),
			makeLexeme("false", lex.LBoolFalse, 6),
			makeLexeme(" ", lex.LWhitespace, 11),
			makeLexeme("TRUE", lex.LBoolTrue, 12),
			makeLexeme(" ", lex.LWhitespace, 16),
			makeLexeme("FALSE", lex.LBoolFalse, 17),
		},
	)
}

func TestMultipleWhitespace(t *testing.T) {
	doTestGetNext(
		t,
		" \t\n",
		[]lex.Lexeme{
			makeLexeme(" \t\n", lex.LWhitespace, 1),
		},
	)
}

func TestIdentifierWithSpecialPrefix(t *testing.T) {
	doTestGetNext(
		t,
		"iffoo",
		[]lex.Lexeme{
			makeLexeme("iffoo", lex.LIdentifier, 1),
		},
	)
}

func TestOperatorBetweenIdentifiers(t *testing.T) {
	doTestGetNext(
		t,
		"foo+bar",
		[]lex.Lexeme{
			makeLexeme("foo", lex.LIdentifier, 1),
			makeLexeme("+", lex.LOperator, 4),
			makeLexeme("bar", lex.LIdentifier, 5),
		},
	)
}

func TestNotSingleOperatorBetweenIdentifiers(t *testing.T) {
	doTestGetNext(
		t,
		"foo++bar",
		[]lex.Lexeme{
			makeLexeme("foo", lex.LIdentifier, 1),
			makeLexeme("++", lex.LOperator, 4),
			makeLexeme("bar", lex.LIdentifier, 6),
		},
	)
}

func TestFunctionDeclaration(t *testing.T) {
	doTestGetNext(
		t,
		"() => {}",
		[]lex.Lexeme{
			makeLexeme("(", lex.LParenOpen, 1),
			makeLexeme(")", lex.LParenClose, 2),
			makeLexeme(" ", lex.LWhitespace, 3),
			makeLexeme("=>", lex.LOperator, 4),
			makeLexeme(" ", lex.LWhitespace, 6),
			makeLexeme("{", lex.LBraceOpen, 7),
			makeLexeme("}", lex.LBraceClose, 8),
		},
	)
}

func TestGte(t *testing.T) {
	doTestGetNext(
		t,
		"3>=4",
		[]lex.Lexeme{
			makeLexeme("3", lex.LNumber, 1),
			makeLexeme(">=", lex.LOperator, 2),
			makeLexeme("4", lex.LNumber, 4),
		},
	)
}

func TestAssignment(t *testing.T) {
	code :=
		`let n : number = 4;
let total : number = 0;`

	doTestGetNext(
		t,
		code,
		[]lex.Lexeme{
			makeLexeme("let", lex.LLet, 1),
			makeLexeme(" ", lex.LWhitespace, 4),
			makeLexeme("n", lex.LIdentifier, 5),
			makeLexeme(" ", lex.LWhitespace, 6),
			makeLexeme(":", lex.LOperator, 7),
			makeLexeme(" ", lex.LWhitespace, 8),
			makeLexeme("number", lex.LIdentifier, 9),
			makeLexeme(" ", lex.LWhitespace, 15),
			makeLexeme("=", lex.LOperator, 16),
			makeLexeme(" ", lex.LWhitespace, 17),
			makeLexeme("4", lex.LNumber, 18),
			makeLexeme(";", lex.LSemiColon, 19),
			makeLexeme("\n", lex.LWhitespace, 20),
			makeLexeme("let", lex.LLet, 21),
			makeLexeme(" ", lex.LWhitespace, 24),
			makeLexeme("total", lex.LIdentifier, 25),
			makeLexeme(" ", lex.LWhitespace, 30),
			makeLexeme(":", lex.LOperator, 31),
			makeLexeme(" ", lex.LWhitespace, 32),
			makeLexeme("number", lex.LIdentifier, 33),
			makeLexeme(" ", lex.LWhitespace, 39),
			makeLexeme("=", lex.LOperator, 40),
			makeLexeme(" ", lex.LWhitespace, 41),
			makeLexeme("0", lex.LNumber, 42),
			makeLexeme(";", lex.LSemiColon, 43),
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

func makeLexeme(value string, lexemeType lex.LexemeType, position int) lex.Lexeme {
	return lex.Lexeme{
		position,
		lexemeType,
		value,
	}
}
