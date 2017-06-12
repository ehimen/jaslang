package run

import (
	"io"

	"github.com/ehimen/jaslang/lex"
	"github.com/ehimen/jaslang/parse"
	"github.com/ehimen/jaslang/runtime"
)

func Interpret(code io.RuneReader, input io.Reader, output io.Writer) error {
	// TODO: should accept error writer here to allow writing
	// TODO: to stdout, rather than having to pass Go errors.
	parser := parse.NewParser(lex.NewJslLexer(code))

	if ast, err := parser.Parse(); err != nil {
		return err
	} else {
		if err := runtime.NewEvaluator(input, output).Evaluate(ast); err != nil {
			return err
		}
	}

	return nil
}
