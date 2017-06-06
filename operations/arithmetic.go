package operations

import (
	"github.com/ehimen/jaslang/functions"
	"github.com/ehimen/jaslang/runtime"
)

type Sum struct{}

func (s Sum) Operator() string {
	return "+"
}

func (s Sum) Precedence() int {
	return 0
}

func (s Sum) Invokable() runtime.Invokable {
	return functions.Add{}
}

type Subtract struct{}

func (s Subtract) Operator() string {
	return "-"
}

func (s Subtract) Precedence() int {
	return 0
}

func (s Subtract) Invokable() runtime.Invokable {
	return functions.Noop{}
}

type Multiply struct{}

func (m Multiply) Operator() string {
	return "*"
}

func (m Multiply) Precedence() int {
	return 1
}

func (m Multiply) Invokable() runtime.Invokable {
	return functions.Noop{}
}
