package operations_test

import (
	"testing"

	"github.com/ehimen/jaslang/operations"
)

type test_operation struct {
	operator   string
	precedence int
}

func (o *test_operation) Operator() string {
	return o.operator
}

func (o *test_operation) Precedence() int {
	return o.precedence
}

func operatorWithPrecedence(operator string, precedence int) operations.Operation {
	return &test_operation{operator: operator, precedence: precedence}
}

func TestPrecedence(t *testing.T) {
	cases := []struct {
		what     operations.Operation
		over     operations.Operation
		expected bool
	}{
		{
			operatorWithPrecedence("+", 0),
			operatorWithPrecedence("-", 0),
			false,
		},
		{
			operatorWithPrecedence("+", 1),
			operatorWithPrecedence("-", 0),
			true,
		},
		{
			operatorWithPrecedence("+", 0),
			operatorWithPrecedence("-", 1),
			false,
		},
		{
			operatorWithPrecedence("+", 1),
			operatorWithPrecedence("-", 1),
			false,
		},
	}

	for _, test := range cases {
		register := operations.NewRegister()
		register.Register(test.what)
		register.Register(test.over)

		actual, err := register.TakesPrecedence(test.what.Operator(), test.over.Operator())

		if err != nil {
			t.Errorf("Unexpected error when determining precedence, %v", err)
		}

		if test.expected && !actual {
			t.Errorf("Expected %s to take precedence over %s, but it didn't", test.what.Operator(), test.over.Operator())
		} else if !test.expected && actual {
			t.Errorf("Expected %s not to take precedence over %s, but it did", test.what.Operator(), test.over.Operator())
		}
	}
}

func TestPrecedenceThrowsWhenInvalidOperator(t *testing.T) {
	register := operations.NewRegister()

	register.Register(operatorWithPrecedence("+", 0))

	_, err := register.TakesPrecedence("+", "-")

	if _, isError := err.(operations.UnknownOperatorError); !isError || err.Error() != "Unknown operator: -" {
		t.Errorf("Expected unknown operator error, but got: %v", err)
	}

	_, err = register.TakesPrecedence("-", "+")

	if _, isError := err.(operations.UnknownOperatorError); !isError || err.Error() != "Unknown operator: -" {
		t.Errorf("Expected unknown operator error, but got: %v", err)
	}
}
