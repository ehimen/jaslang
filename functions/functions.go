package functions

import (
	"errors"

	"github.com/ehimen/jaslang/runtime"
)

// Contains native functions.
// This file probably splits out in to
// groups of functions.

type Println struct {
}

func (p Println) String() string {
	return "println() <native>"
}

func (p Println) Invoke(context *runtime.Context, args []runtime.Value) (error, runtime.Value) {
	for _, arg := range args {
		context.Output.Write([]byte(arg.String()))
	}

	return nil, runtime.Void{}
}

type Add struct {
}

func (p Add) String() string {
	return "+ <native>"
}

func (a Add) Invoke(context *runtime.Context, args []runtime.Value) (error, runtime.Value) {
	if one, isNumber := args[0].(runtime.Number); isNumber {
		if two, isNumber := args[1].(runtime.Number); isNumber {
			return nil, runtime.Number{Value: one.Value + two.Value}
		}
	}

	return errors.New("Invalid operands. + requires two numbers"), nil
}

type Noop struct {
}

func (n Noop) String() string {
	return "noop <native>"
}

func (n Noop) Invoke(context *runtime.Context, args []runtime.Value) (error, runtime.Value) {
	return nil, runtime.Void{}
}
