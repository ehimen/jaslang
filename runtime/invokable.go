package runtime

type Invokable interface {
	Value
	Invoke(context *Context, args []Value) (error, Value)
}
