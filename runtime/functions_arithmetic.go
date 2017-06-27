package runtime

import "errors"

type AddNumbers struct {
}

func (a AddNumbers) String() string {
	return "addition() <native>"
}

func (a AddNumbers) Type() Type {
	return TypeInvokable
}

func (a AddNumbers) Invoke(context *Context, args []Value) (error, Value) {
	if one, isNumber := args[0].(Number); isNumber {
		if two, isNumber := args[1].(Number); isNumber {
			return nil, Number{Value: one.Value + two.Value}
		}
	}

	return errors.New("Invalid operands. Number addition requires two numbers"), nil
}

type SubtractNumbers struct {
}

func (a SubtractNumbers) String() string {
	return "subtraction() <native>"
}

func (a SubtractNumbers) Type() Type {
	return TypeInvokable
}

func (a SubtractNumbers) Invoke(context *Context, args []Value) (error, Value) {
	if one, isNumber := args[0].(Number); isNumber {
		if two, isNumber := args[1].(Number); isNumber {
			return nil, Number{Value: one.Value - two.Value}
		}
	}

	return errors.New("Invalid operands. Subtraction requires two numbers"), nil
}

type MultiplyNumbers struct {
}

func (m MultiplyNumbers) String() string {
	return "multiply() <native>"
}

func (m MultiplyNumbers) Type() Type {
	return TypeInvokable
}

func (m MultiplyNumbers) Invoke(context *Context, args []Value) (error, Value) {
	if one, isNumber := args[0].(Number); isNumber {
		if two, isNumber := args[1].(Number); isNumber {
			return nil, Number{Value: one.Value * two.Value}
		}
	}

	return errors.New("Invalid operands. Multiplication requires two numbers"), nil
}

type DivideNumbers struct {
}

func (d DivideNumbers) String() string {
	return "divide() <native>"
}

func (d DivideNumbers) Type() Type {
	return TypeInvokable
}

func (d DivideNumbers) Invoke(context *Context, args []Value) (error, Value) {
	if one, isNumber := args[0].(Number); isNumber {
		if two, isNumber := args[1].(Number); isNumber {
			return nil, Number{Value: one.Value / two.Value}
		}
	}

	return errors.New("Invalid operands. Division requires two numbers"), nil
}
