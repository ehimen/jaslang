package runtime

import (
	"io"

	"errors"
	"fmt"

	"github.com/ehimen/jaslang/parse"
)

type Evaluator interface {
	Evaluate(parse.Node) error
}

type evaluator struct {
	context *Context
}

func NewEvaluator(input io.Reader, output io.Writer, error io.Writer) Evaluator {
	table := NewTable()

	table.AddFunction("println", Println{})
	table.AddOperator("+", Types([]Type{TypeNumber, TypeNumber}), AddNumbers{})
	table.AddOperator("-", Types([]Type{TypeNumber, TypeNumber}), SubtractNumbers{})
	table.AddOperator("+", Types([]Type{TypeString, TypeString}), StringConcatenation{})

	return &evaluator{context: &Context{Table: table, Input: input, Output: output, Error: error}}
}

func (e *evaluator) Evaluate(node parse.Node) error {

	if err, _ := e.evaluate(node); err != nil {
		return err
	}

	return nil
}

func (e *evaluator) evaluate(node parse.Node) (error, Value) {
	args := []Value{}

	if parent, isParent := node.(parse.ContainsChildren); isParent {
		for _, child := range parent.Children() {
			// TODO: not recursion to avoid stack overflows.
			if err, arg := e.evaluate(child); err != nil {
				return err, nil
			} else {
				args = append(args, arg)
			}
		}
	} else if root, isRoot := node.(parse.RootNode); isRoot {
		for _, child := range root.Statements {
			// TODO: not recursion to avoid stack overflows.
			if err, arg := e.evaluate(child); err != nil {
				return err, nil
			} else {
				args = append(args, arg)
			}
		}
	}

	if str, isStr := node.(*parse.String); isStr {
		return nil, String{Value: str.Value}
	}

	if num, isNum := node.(*parse.Number); isNum {
		return nil, Number{Value: num.Value}
	}

	if boolean, isBool := node.(*parse.Boolean); isBool {
		return nil, Boolean{Value: boolean.Value}
	}

	if fn, isFn := node.(*parse.FunctionCall); isFn {
		return e.evaluateFunctionCall(fn, args)
	}

	if operator, isOperator := node.(*parse.Operator); isOperator {
		return e.evaluateOperator(operator, args)
	}

	// Nothing to do with statements/root as these are AST constructs (for now).
	if _, isStmt := node.(*parse.Statement); isStmt {
		return nil, nil
	}

	if _, isRoot := node.(parse.RootNode); isRoot {
		return nil, nil
	}

	return errors.New(fmt.Sprintf("Handling for %#v not yet implemented.", node)), nil
}

func (e *evaluator) evaluateFunctionCall(fn *parse.FunctionCall, args []Value) (error, Value) {
	if invokable, err := e.context.Table.Invokable(fn.Identifier.Identifier); err != nil {
		return err, nil
	} else {
		invokable.Invoke(e.context, args)
	}

	return nil, nil
}

func (e *evaluator) evaluateOperator(operator *parse.Operator, args []Value) (error, Value) {
	operands := Types([]Type{})

	for _, arg := range args {
		operands = append(operands, arg.Type())
	}

	if invokable, err := e.context.Table.Operator(operator.Operator, operands); err != nil {
		return err, nil
	} else {
		return invokable.Invoke(e.context, args)
	}
}
