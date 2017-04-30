package operations

type Sum struct{}

func (s Sum) Operator() string {
	return "+"
}

func (s Sum) Precedence() int {
	return 0
}

type Multiply struct{}

func (m Multiply) Operator() string {
	return "*"
}

func (m Multiply) Precedence() int {
	return 1
}
