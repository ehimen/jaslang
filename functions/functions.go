package functions

import (
	"fmt"

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
	// TODO: For some reason, context's input & output are nil when they get here,
	// TODO: They look good when this is called, just borked when we arrive =/
	fmt.Printf("%+v\n", context)
	context.Writer.Write([]byte("foo"))
	//for _ := range args {
	//}
}
