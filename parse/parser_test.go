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

	expected := parse.RootNode{
		Statements: []*parse.Statement{
			{
				ParentNode: parse.ParentNode{
					Children: []parse.Node{
						&parse.FunctionCall{
							Identifier: "print",
							ParentNode: parse.ParentNode{
								Children: []parse.Node{
									&parse.StringLiteral{
										Value: "Hello, world!",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, parser.Parse())
}

func getParser(lexemes []lex.Lexeme) parse.Parser {
	return parse.NewParser(testutil.NewSimpleLexer(lexemes))
}
