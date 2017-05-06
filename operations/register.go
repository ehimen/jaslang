package operations

import (
	"fmt"
)

// Defines an operation (e.g. +, * etc)
// For now, this simply defines precedence,
// but we intend that implementations of this
// can define the operation itself, and using
// reflection we can resolve an operator
// to a function that will actually perform the
// operation. Reflection _may_ allow us to have
// different signatures for the same operator,
// (e.g. "foo" + "bar" vs 1 + 2), but this will
// need some more thinking, especially around
// any cross over with functions. To proceed here,
// we need to flesh out the evaluator more and consider
// how values are represented and passed to functions.
type Operation interface {
	Operator() string
	Precedence() int
}

type Register struct {
	operations map[string]Operation
}

type UnknownOperatorError struct {
	operator string
}

func (err UnknownOperatorError) Error() string {
	return fmt.Sprintf("Unknown operator: %s", err.operator)
}

func NewRegister() *Register {
	return &Register{operations: make(map[string]Operation)}
}

func (r *Register) Register(operation Operation) error {
	r.operations[operation.Operator()] = operation

	return nil
}

// Returns true if what is higher precedence than over,
// i.e. what should be evaluated before over.
// If either operator is unknown to this register,
// this will return false.
func (r *Register) TakesPrecedence(what string, over string) (bool, error) {
	if whatOperator, exists := r.operations[what]; !exists {
		return false, UnknownOperatorError{operator: what}
	} else if overOperator, exists := r.operations[over]; !exists {
		return false, UnknownOperatorError{operator: over}
	} else {
		return whatOperator.Precedence() > overOperator.Precedence(), nil
	}
}
