package runtime

type Println struct {
}

func (p Println) String() string {
	return "println() <native>"
}

func (p Println) Type() Type {
	return TypeInvokable
}

func (p Println) Invoke(context *Context, args []Value) (error, Value) {
	for _, arg := range args {
		context.Output.Write([]byte(arg.String()))
	}

	return nil, Void{}
}
