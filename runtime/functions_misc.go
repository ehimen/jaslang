package runtime

type Noop struct {
}

func (n Noop) Type() Type {
	return TypeInvokable
}

func (n Noop) String() string {
	return "noop <native>"
}

func (n Noop) Invoke(context *Context, args []Value) (error, Value) {
	return nil, Void{}
}
