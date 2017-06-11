package runtime

import "errors"

type StringConcatenation struct {
}

func (p StringConcatenation) String() string {
	return "StringConcatenation() <native>"
}

func (p StringConcatenation) Type() Type {
	return TypeInvokable
}

func (p StringConcatenation) Invoke(context *Context, args []Value) (error, Value) {
	if one, isString := args[0].(String); isString {
		if two, isString := args[1].(String); isString {
			return nil, String{Value: one.Value + two.Value}
		}
	}

	return errors.New("Invalid operands. String concatenation requires two strings"), nil
}
