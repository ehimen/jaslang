package parse

import (
	"fmt"
)

// Defines an operation (e.g. +, * etc)
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
