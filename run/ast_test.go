package run_test

import (
	"encoding/json"
	"testing"

	"strings"

	"github.com/ehimen/jaslang/lex"
	"github.com/ehimen/jaslang/parse"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name     string
	input    string
	expected string
}{
	{
		name:  "function call",
		input: "println(1 + 2 - 3);",
		expected: `
[
	{
		"Type": "statement",
		"Children": [
			{
				"Type": "function",
				"Identifier": "println",
				"Children": [
					{
						"Type": "operator",
						"Operator": "-",
						"Children": [
							{
								"Type": "operator",
								"Operator": "+",
								"Children": [
									{
										"Type": "number",
										"Value": 1
									},
									{
										"Type": "number",
										"Value": 2
									}
								]
							},
							{
								"Type": "number",
								"Value": 3
							}
						]
					}
				]
			}
		]
	}
]
`,
	},
	{
		name:  "declaration + assignment",
		input: "let b boolean = true;",
		expected: `
[
	{
		"Type": "statement",
		"Children": [
			{
				"Type": "declaration",
				"Identifier": "b",
				"ValueType": "boolean",
				"Children": [
					{
						"Type": "boolean",
						"Value": true
					}
				]
			}
		]
	}
]
`,
	},
}

func TestAst(t *testing.T) {

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				parser := parse.NewParser(lex.NewJslLexer(strings.NewReader(test.input)))

				if ast, err := parser.Parse(); err != nil {
					t.Fatalf("Could not parse input: %s", err)
				} else {
					if actual, err := json.Marshal(ast); err != nil {
						t.Fatalf("Error encoding AST: %s", err)
					} else {
						assert.JSONEq(t, string(test.expected), string(actual))
					}
				}
			},
		)
	}
}
