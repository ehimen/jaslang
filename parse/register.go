package parse

import (
	"fmt"
)

type Register struct {
	operations map[string]int
}

type UnknownOperatorError struct {
	operator string
}

func (err UnknownOperatorError) Error() string {
	return fmt.Sprintf("Unknown operator: %s", err.operator)
}

func NewRegister() *Register {
	return &Register{operations: make(map[string]int)}
}

func (r *Register) Register(operator string, precedence int) error {
	r.operations[operator] = precedence

	return nil
}

// Returns true if what is higher precedence than over,
// i.e. what should be evaluated before over.
// If either operator is unknown to this register,
// this will return false.
func (r *Register) TakesPrecedence(what string, over string) (bool, error) {
	if whatPrecedence, exists := r.operations[what]; !exists {
		return false, UnknownOperatorError{operator: what}
	} else if overPrecedence, exists := r.operations[over]; !exists {
		return false, UnknownOperatorError{operator: over}
	} else {
		return whatPrecedence > overPrecedence, nil
	}
}
