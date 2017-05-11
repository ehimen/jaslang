package run

import (
	"io"

	"github.com/ehimen/jaslang/evaluation"
	"github.com/ehimen/jaslang/lex"
	"github.com/ehimen/jaslang/parse"
)

func Interpret(code io.RuneReader, input io.Reader, output io.Writer) error {
	parser := parse.NewParser(lex.NewJslLexer(code))

	if ast, err := parser.Parse(); err != nil {
		return err
	} else {
		if err := evaluation.NewEvaluator(input, output).Evaluate(ast); err != nil {
			return err
		}
	}

	return nil
}
