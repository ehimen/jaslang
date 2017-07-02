package parse_test

import (
	"testing"

	"github.com/ehimen/jaslang/lex"
	"github.com/ehimen/jaslang/parse"
	"github.com/ehimen/jaslang/testutil"

	"github.com/stretchr/testify/assert"
)

func TestSimpleFunctionCall(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("print", lex.LIdentifier, 1, 1),
		testutil.MakeLexeme("(", lex.LParenOpen, 2, 1),
		testutil.MakeLexeme("Hello, world!", lex.LQuoted, 3, 1),
		testutil.MakeLexeme(")", lex.LParenClose, 4, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 5, 1),
	})

	expected := expectStatements(
		parse.NewStatement(
			1,
			1,
			parse.NewFunctionCall(
				"print",
				1,
				1,
				parse.NewString("Hello, world!", 1, 3),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestTwoLiterals(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("1.34", lex.LNumber, 1, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 2, 1),
		testutil.MakeLexeme(" ", lex.LWhitespace, 3, 1),
		testutil.MakeLexeme("3.42", lex.LNumber, 4, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 5, 1),
	})

	expected := expectStatements(
		parse.NewStatement(1, 1, parse.NewNumber(float64(1.34), 1, 1)),
		parse.NewStatement(1, 4, parse.NewNumber(float64(3.42), 1, 4)),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestInvalidNumberSyntax(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("1.3.2.2.422", lex.LNumber, 1, 1),
	})

	_, err := parser.Parse()

	if _, isInvalidNumber := err.(parse.InvalidNumberError); !isInvalidNumber {
		t.Fatalf("Expected Parse() to fail on invalid number, but got: %s", err)
	}
}

func TestIncompleteInput(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("true", lex.LBoolTrue, 1, 1),
	})

	if _, err := parser.Parse(); err != parse.UnterminatedStatement {
		t.Fatalf("Expected unterminated statement error, but got: %v", err)
	}
}

func TestTrueFalse(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("true", lex.LBoolTrue, 1, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 2, 1),
		testutil.MakeLexeme(" ", lex.LWhitespace, 3, 1),
		testutil.MakeLexeme("false", lex.LBoolFalse, 4, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 5, 1),
	})

	expected := expectStatements(
		parse.NewStatement(1, 1, parse.NewBoolean(true, 1, 1)),
		parse.NewStatement(1, 4, parse.NewBoolean(false, 1, 4)),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestOperator(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("a", lex.LIdentifier, 1, 1),
		testutil.MakeLexeme("+", lex.LOperator, 2, 1),
		testutil.MakeLexeme("b", lex.LIdentifier, 3, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 4, 1),
	})

	expected := expectStatements(
		parse.NewStatement(
			1,
			1,
			parse.NewOperator(
				"+",
				1,
				2,
				parse.NewIdentifier("a", 1, 1),
				parse.NewIdentifier("b", 1, 3),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestMultipleOperator(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("a", lex.LIdentifier, 1, 1),
		testutil.MakeLexeme("+", lex.LOperator, 2, 1),
		testutil.MakeLexeme("b", lex.LIdentifier, 3, 1),
		testutil.MakeLexeme("+", lex.LOperator, 4, 1),
		testutil.MakeLexeme("c", lex.LIdentifier, 5, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 6, 1),
	})

	expected := expectStatements(
		parse.NewStatement(
			1,
			1,
			parse.NewOperator(
				"+",
				1,
				4,
				parse.NewOperator(
					"+",
					1,
					2,
					parse.NewIdentifier("a", 1, 1),
					parse.NewIdentifier("b", 1, 3),
				),
				parse.NewIdentifier("c", 1, 5),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestLet(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1, 1),
		testutil.MakeLexeme("foo", lex.LIdentifier, 2, 1),
		testutil.MakeLexeme("number", lex.LIdentifier, 3, 1),
		testutil.MakeLexeme("=", lex.LEquals, 4, 1),
		testutil.MakeLexeme("1", lex.LNumber, 5, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 6, 1),
	})

	expected := expectStatements(
		parse.NewStatement(
			1,
			1,
			parse.NewDeclaration(
				*parse.NewIdentifier("foo", 1, 2),
				*parse.NewIdentifier("number", 1, 3),
				1,
				1,
				parse.NewNumber(1, 1, 5),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestLetWithExpression(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1, 1),
		testutil.MakeLexeme("foo", lex.LIdentifier, 2, 1),
		testutil.MakeLexeme("number", lex.LIdentifier, 3, 1),
		testutil.MakeLexeme("=", lex.LEquals, 4, 1),
		testutil.MakeLexeme("1", lex.LNumber, 5, 1),
		testutil.MakeLexeme("+", lex.LOperator, 6, 1),
		testutil.MakeLexeme("2", lex.LNumber, 7, 1),
		testutil.MakeLexeme("-", lex.LOperator, 8, 1),
		testutil.MakeLexeme("3", lex.LNumber, 9, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 10, 1),
	})

	expected := expectStatements(
		parse.NewStatement(
			1,
			1,
			parse.NewDeclaration(
				*parse.NewIdentifier("foo", 1, 2),
				*parse.NewIdentifier("number", 1, 3),
				1,
				1,
				parse.NewOperator(
					"-",
					1,
					8,
					parse.NewOperator(
						"+",
						1,
						6,
						parse.NewNumber(1, 1, 5),
						parse.NewNumber(2, 1, 7),
					),
					parse.NewNumber(3, 1, 9),
				),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestOperatorPrecedence(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("1", lex.LNumber, 1, 1),
		testutil.MakeLexeme("*", lex.LOperator, 2, 1),
		testutil.MakeLexeme("2", lex.LNumber, 3, 1),
		testutil.MakeLexeme("+", lex.LOperator, 4, 1),
		testutil.MakeLexeme("3", lex.LNumber, 5, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 6, 1),
	})

	expected := expectStatements(
		parse.NewStatement(
			1,
			1,
			parse.NewOperator(
				"+",
				1,
				4,
				parse.NewOperator(
					"*",
					1,
					2,
					parse.NewNumber(1, 1, 1),
					parse.NewNumber(2, 1, 3),
				),
				parse.NewNumber(3, 1, 5),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestOperatorAndFunction(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("a", lex.LIdentifier, 1, 1),
		testutil.MakeLexeme("+", lex.LOperator, 2, 1),
		testutil.MakeLexeme("foo", lex.LIdentifier, 3, 1),
		testutil.MakeLexeme("(", lex.LParenOpen, 4, 1),
		testutil.MakeLexeme("b", lex.LIdentifier, 5, 1),
		testutil.MakeLexeme(")", lex.LParenClose, 6, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 7, 1),
	})

	expected := expectStatements(
		parse.NewStatement(
			1,
			1,
			parse.NewOperator(
				"+",
				1,
				2,
				parse.NewIdentifier("a", 1, 1),
				parse.NewFunctionCall(
					"foo",
					1,
					3,
					parse.NewIdentifier("b", 1, 5),
				),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestInvalidLetAssigned(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1, 1),
		testutil.MakeLexeme("=", lex.LEquals, 2, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 3, 1),
	})

	_, err := parser.Parse()

	if unexpectedToken, isUnexpectedToken := err.(parse.UnexpectedTokenError); !isUnexpectedToken {
		t.Fatalf("Expected unexpected token error, but got: %v", err)
	} else {
		assert.Equal(t, "Unexpected token \"=\" (position 2, line 1)", unexpectedToken.Error())
	}
}

func TestInvalidLetWithoutType(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1, 1),
		testutil.MakeLexeme("foo", lex.LIdentifier, 2, 1),
		testutil.MakeLexeme("=", lex.LEquals, 3, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 4, 1),
	})

	_, err := parser.Parse()

	if unexpectedToken, isUnexpectedToken := err.(parse.UnexpectedTokenError); !isUnexpectedToken {
		t.Fatalf("Expected unexpected token error, but got: %v", err)
	} else {
		assert.Equal(t, "Unexpected token \"=\" (position 3, line 1)", unexpectedToken.Error())
	}
}

func TestInvalidNestedLet(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1, 1),
		testutil.MakeLexeme("foo", lex.LIdentifier, 2, 1),
		testutil.MakeLexeme("string", lex.LIdentifier, 3, 1),
		testutil.MakeLexeme("=", lex.LEquals, 4, 1),
		testutil.MakeLexeme("let", lex.LLet, 5, 1),
	})

	_, err := parser.Parse()

	if unexpectedToken, isUnexpectedToken := err.(parse.UnexpectedTokenError); !isUnexpectedToken {
		t.Fatalf("Expected unexpected token error, but got: %v", err)
	} else {
		assert.Equal(t, "Unexpected token \"let\" (position 5, line 1)", unexpectedToken.Error())
	}
}

func getParser(lexemes []lex.Lexeme) parse.Parser {
	return parse.NewParser(testutil.NewSimpleLexer(lexemes))
}

func expectStatements(statements ...*parse.Statement) parse.RootNode {
	return parse.RootNode{Statements: statements}
}

func testParse(p parse.Parser, t *testing.T) parse.RootNode {
	node, err := p.Parse()

	if err != nil {
		if invalidTokenErr, isInvalidToken := err.(parse.UnexpectedTokenError); isInvalidToken {
			t.Fatalf("%v, debug: %s", invalidTokenErr, invalidTokenErr.Debug)
		} else {
			t.Fatalf("Unexpected Parse() error: %v", err)
		}
	}

	return node
}
