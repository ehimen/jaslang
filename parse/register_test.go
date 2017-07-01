package parse_test

import (
	"testing"

	"github.com/ehimen/jaslang/parse"
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

func operatorWithPrecedence(operator string, precedence int) *test_operation {
	return &test_operation{operator: operator, precedence: precedence}
}

func TestPrecedence(t *testing.T) {
	cases := []struct {
		what     *test_operation
		over     *test_operation
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
		register := parse.NewRegister()
		register.Register(test.what.Operator(), test.what.Precedence())
		register.Register(test.over.Operator(), test.over.Precedence())

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
	register := parse.NewRegister()

	register.Register("+", 0)

	_, err := register.TakesPrecedence("+", "-")

	if _, isError := err.(parse.UnknownOperatorError); !isError || err.Error() != "Unknown operator: -" {
		t.Errorf("Expected unknown operator error, but got: %v", err)
	}

	_, err = register.TakesPrecedence("-", "+")

	if _, isError := err.(parse.UnknownOperatorError); !isError || err.Error() != "Unknown operator: -" {
		t.Errorf("Expected unknown operator error, but got: %v", err)
	}
}
