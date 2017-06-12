package parse

type Sum struct{}

func (s Sum) Operator() string {
	return "+"
}

func (s Sum) Precedence() int {
	return 0
}

type Subtract struct{}

func (s Subtract) Operator() string {
	return "-"
}

func (s Subtract) Precedence() int {
	return 0
}

type Multiply struct{}

func (m Multiply) Operator() string {
	return "*"
}

func (m Multiply) Precedence() int {
	return 1
}
