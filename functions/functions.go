package functions

import (
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

func (p Println) Invoke(context *runtime.Context, args []runtime.Value) {
	for _, arg := range args {
		context.Output.Write([]byte(arg.String()))
	}
}
