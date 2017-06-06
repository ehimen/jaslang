package evaluation

import (
	"io"

	"github.com/ehimen/jaslang/functions"
	"github.com/ehimen/jaslang/parse"
	"github.com/ehimen/jaslang/runtime"
)

type Evaluator interface {
	Evaluate(parse.Node) error
}

type evaluator struct {
	context *runtime.Context
}

func NewEvaluator(input io.Reader, output io.Writer) Evaluator {
	table := runtime.NewTable()

	table.AddFunction("println", functions.Println{})

	return &evaluator{context: &runtime.Context{Table: table, Input: input, Output: output}}
}

func (e *evaluator) Evaluate(node parse.Node) error {

	if _, err := e.evaluate(node); err != nil {
		return err
	}

	return nil
}

func (e *evaluator) evaluate(node parse.Node) (runtime.Value, error) {
	args := []runtime.Value{}

	if parent, isParent := node.(parse.ContainsChildren); isParent {
		for _, child := range parent.Children() {
			// TODO: not recursion to avoid stack overflows.
			if arg, err := e.evaluate(child); err != nil {
				return nil, err
			} else {
				args = append(args, arg)
			}
		}
	} else if root, isRoot := node.(parse.RootNode); isRoot {
		for _, child := range root.Statements {
			// TODO: not recursion to avoid stack overflows.
			if arg, err := e.evaluate(child); err != nil {
				return nil, err
			} else {
				args = append(args, arg)
			}
		}
	}

	if str, isStr := node.(*parse.String); isStr {
		return runtime.String{Value: str.Value}, nil
	}

	if num, isNum := node.(*parse.Number); isNum {
		return runtime.Number{Value: num.Value}, nil
	}

	if boolean, isBool := node.(*parse.Boolean); isBool {
		return runtime.Boolean{Value: boolean.Value}, nil
	}

	if fn, isFn := node.(*parse.FunctionCall); isFn {
		return e.evaluateFunctionCall(fn, args)
	}

	return nil, nil
}

func (e *evaluator) evaluateFunctionCall(fn *parse.FunctionCall, args []runtime.Value) (runtime.Value, error) {
	if invokable, err := e.context.Table.Invokable(fn.Identifier.Identifier); err != nil {
		return nil, err
	} else {
		invokable.Invoke(e.context, args)
	}

	return nil, nil
}

func (e *evaluator) evaluateOperator(fn *parse.Operator, args []runtime.Value) (runtime.Value, error) {
	e.context.
}
