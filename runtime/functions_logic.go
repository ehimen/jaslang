package runtime

import "errors"

type LogicAnd struct {
}

func (l LogicAnd) String() string {
	return "&& <native>"
}

func (l LogicAnd) Type() Type {
	return TypeInvokable
}

func (l LogicAnd) Invoke(context *Context, args []Value) (error, Value) {
	if one, isBoolean := args[0].(Boolean); isBoolean {
		if two, isBoolean := args[1].(Boolean); isBoolean {
			return nil, Boolean{Value: one.Value && two.Value}
		}
	}

	return errors.New("Invalid operands. Logic AND requires two booleans"), nil
}

type LogicOr struct {
}

func (l LogicOr) String() string {
	return "|| <native>"
}

func (l LogicOr) Type() Type {
	return TypeInvokable
}

func (l LogicOr) Invoke(context *Context, args []Value) (error, Value) {
	if one, isBoolean := args[0].(Boolean); isBoolean {
		if two, isBoolean := args[1].(Boolean); isBoolean {
			return nil, Boolean{Value: one.Value || two.Value}
		}
	}

	return errors.New("Invalid operands. Logic OR requires two booleans"), nil
}

type Equality struct {
}

func (l Equality) String() string {
	return "== <native>"
}

func (l Equality) Type() Type {
	return TypeInvokable
}

func (l Equality) Invoke(context *Context, args []Value) (error, Value) {
	if one, isNumber := args[0].(Number); isNumber {
		if two, isNumber := args[1].(Number); isNumber {
			return nil, Boolean{Value: one.Value == two.Value}
		}
	}

	// TODO: we should be able to perform equality for any kind of operator.
	return errors.New("Invalid operands. Equality requires two numbers"), nil
}

type LessThan struct {
}

func (l LessThan) String() string {
	return "< <native>"
}

func (l LessThan) Type() Type {
	return TypeInvokable
}

func (l LessThan) Invoke(context *Context, args []Value) (error, Value) {
	if one, isNumber := args[0].(Number); isNumber {
		if two, isNumber := args[1].(Number); isNumber {
			return nil, Boolean{Value: one.Value < two.Value}
		}
	}

	return errors.New("Invalid operands. Less than comparison requires two numbers"), nil
}

type GreaterThan struct {
}

func (l GreaterThan) String() string {
	return "< <native>"
}

func (l GreaterThan) Type() Type {
	return TypeInvokable
}

func (l GreaterThan) Invoke(context *Context, args []Value) (error, Value) {
	if one, isNumber := args[0].(Number); isNumber {
		if two, isNumber := args[1].(Number); isNumber {
			return nil, Boolean{Value: one.Value < two.Value}
		}
	}

	return errors.New("Invalid operands. Greater than comparison requires two numbers"), nil
}
