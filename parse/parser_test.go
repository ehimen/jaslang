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
		testutil.MakeLexeme("print", lex.LIdentifier, 1),
		testutil.MakeLexeme("(", lex.LParenOpen, 2),
		testutil.MakeLexeme("Hello, world!", lex.LQuoted, 3),
		testutil.MakeLexeme(")", lex.LParenClose, 4),
		testutil.MakeLexeme(";", lex.LSemiColon, 5),
	})

	expected := expectStatements(
		parse.NewStatement(
			parse.NewFunctionCall(
				"print",
				parse.NewString("Hello, world!"),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestTwoLiterals(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("1.34", lex.LNumber, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 2),
		testutil.MakeLexeme(" ", lex.LWhitespace, 3),
		testutil.MakeLexeme("3.42", lex.LNumber, 4),
		testutil.MakeLexeme(";", lex.LSemiColon, 5),
	})

	expected := expectStatements(
		parse.NewStatement(parse.NewNumber(float64(1.34))),
		parse.NewStatement(parse.NewNumber(float64(3.42))),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestInvalidNumberSyntax(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("1.3.2.2.422", lex.LNumber, 1),
	})

	_, err := parser.Parse()

	if _, isInvalidNumber := err.(parse.InvalidNumberError); !isInvalidNumber {
		t.Fatalf("Expected Parse() to fail on invalid number, but got: %s", err)
	}
}

func TestIncompleteInput(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("true", lex.LBoolTrue, 1),
	})

	if _, err := parser.Parse(); err != parse.UnterminatedStatement {
		t.Fatalf("Expected unterminated statement error, but got: %v", err)
	}
}

func TestTrueFalse(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("true", lex.LBoolTrue, 1),
		testutil.MakeLexeme(";", lex.LSemiColon, 2),
		testutil.MakeLexeme(" ", lex.LWhitespace, 3),
		testutil.MakeLexeme("false", lex.LBoolFalse, 4),
		testutil.MakeLexeme(";", lex.LSemiColon, 5),
	})

	expected := expectStatements(
		parse.NewStatement(parse.NewBoolean(true)),
		parse.NewStatement(parse.NewBoolean(false)),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestOperator(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("a", lex.LIdentifier, 1),
		testutil.MakeLexeme("+", lex.LOperator, 2),
		testutil.MakeLexeme("b", lex.LIdentifier, 3),
		testutil.MakeLexeme(";", lex.LSemiColon, 4),
	})

	expected := expectStatements(
		parse.NewStatement(
			parse.NewOperator(
				"+",
				parse.NewIdentifier("a"),
				parse.NewIdentifier("b"),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestMultipleOperator(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("a", lex.LIdentifier, 1),
		testutil.MakeLexeme("+", lex.LOperator, 2),
		testutil.MakeLexeme("b", lex.LIdentifier, 3),
		testutil.MakeLexeme("+", lex.LOperator, 4),
		testutil.MakeLexeme("c", lex.LIdentifier, 5),
		testutil.MakeLexeme(";", lex.LSemiColon, 6),
	})

	expected := expectStatements(
		parse.NewStatement(
			parse.NewOperator(
				"+",
				parse.NewOperator(
					"+",
					parse.NewIdentifier("a"),
					parse.NewIdentifier("b"),
				),
				parse.NewIdentifier("c"),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestLet(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1),
		testutil.MakeLexeme("foo", lex.LIdentifier, 2),
		testutil.MakeLexeme("number", lex.LIdentifier, 3),
		testutil.MakeLexeme("=", lex.LEquals, 4),
		testutil.MakeLexeme("1", lex.LNumber, 5),
		testutil.MakeLexeme(";", lex.LSemiColon, 6),
	})

	expected := expectStatements(
		parse.NewStatement(
			&parse.Let{
				Identifier: parse.NewIdentifier("foo"),
				Type:       parse.NewIdentifier("number"),
				Children: []parse.Node{
					parse.NewNumber(1),
				},
			},
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestLetWithExpression(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1),
		testutil.MakeLexeme("foo", lex.LIdentifier, 2),
		testutil.MakeLexeme("number", lex.LIdentifier, 3),
		testutil.MakeLexeme("=", lex.LEquals, 4),
		testutil.MakeLexeme("1", lex.LNumber, 5),
		testutil.MakeLexeme("+", lex.LOperator, 6),
		testutil.MakeLexeme("2", lex.LNumber, 7),
		testutil.MakeLexeme("-", lex.LOperator, 8),
		testutil.MakeLexeme("3", lex.LNumber, 9),
		testutil.MakeLexeme(";", lex.LSemiColon, 10),
	})

	expected := expectStatements(
		parse.NewStatement(
			&parse.Let{
				Identifier: parse.NewIdentifier("foo"),
				Type:       parse.NewIdentifier("number"),
				Children: []parse.Node{
					parse.NewOperator(
						"-",
						parse.NewOperator(
							"+",
							parse.NewNumber(1),
							parse.NewNumber(2),
						),
						parse.NewNumber(3),
					),
				},
			},
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestOperatorPrecedence(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("1", lex.LNumber, 1),
		testutil.MakeLexeme("*", lex.LOperator, 2),
		testutil.MakeLexeme("2", lex.LNumber, 3),
		testutil.MakeLexeme("+", lex.LOperator, 4),
		testutil.MakeLexeme("3", lex.LNumber, 5),
		testutil.MakeLexeme(";", lex.LSemiColon, 6),
	})

	expected := expectStatements(
		parse.NewStatement(
			parse.NewOperator(
				"+",
				parse.NewOperator(
					"*",
					parse.NewNumber(1),
					parse.NewNumber(2),
				),
				parse.NewNumber(3),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestOperatorAndFunction(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("a", lex.LIdentifier, 1),
		testutil.MakeLexeme("+", lex.LOperator, 2),
		testutil.MakeLexeme("foo", lex.LIdentifier, 3),
		testutil.MakeLexeme("(", lex.LParenOpen, 4),
		testutil.MakeLexeme("b", lex.LIdentifier, 5),
		testutil.MakeLexeme(")", lex.LParenClose, 6),
		testutil.MakeLexeme(";", lex.LSemiColon, 7),
	})

	expected := expectStatements(
		parse.NewStatement(
			parse.NewOperator(
				"+",
				parse.NewIdentifier("a"),
				parse.NewFunctionCall(
					"foo",
					parse.NewIdentifier("b"),
				),
			),
		),
	)

	assert.Equal(t, expected, testParse(parser, t))
}

func TestInvalidLetAssigned(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1),
		testutil.MakeLexeme("=", lex.LEquals, 2),
		testutil.MakeLexeme(";", lex.LSemiColon, 3),
	})

	_, err := parser.Parse()

	if unexpectedToken, isUnexpectedToken := err.(parse.UnexpectedTokenError); !isUnexpectedToken {
		t.Fatalf("Expected unexpected token error, but got: %v", err)
	} else {
		assert.Equal(t, "Unexpected token \"=\" at position 2", unexpectedToken.Error())
	}
}

func TestInvalidLetWithoutType(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1),
		testutil.MakeLexeme("foo", lex.LIdentifier, 2),
		testutil.MakeLexeme("=", lex.LEquals, 3),
		testutil.MakeLexeme(";", lex.LSemiColon, 4),
	})

	_, err := parser.Parse()

	if unexpectedToken, isUnexpectedToken := err.(parse.UnexpectedTokenError); !isUnexpectedToken {
		t.Fatalf("Expected unexpected token error, but got: %v", err)
	} else {
		assert.Equal(t, "Unexpected token \"=\" at position 3", unexpectedToken.Error())
	}
}

func TestInvalidNestedLet(t *testing.T) {
	parser := getParser([]lex.Lexeme{
		testutil.MakeLexeme("let", lex.LLet, 1),
		testutil.MakeLexeme("foo", lex.LIdentifier, 2),
		testutil.MakeLexeme("string", lex.LIdentifier, 3),
		testutil.MakeLexeme("=", lex.LEquals, 4),
		testutil.MakeLexeme("let", lex.LLet, 5),
	})

	_, err := parser.Parse()

	if unexpectedToken, isUnexpectedToken := err.(parse.UnexpectedTokenError); !isUnexpectedToken {
		t.Fatalf("Expected unexpected token error, but got: %v", err)
	} else {
		assert.Equal(t, "Unexpected token \"let\" at position 5", unexpectedToken.Error())
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
