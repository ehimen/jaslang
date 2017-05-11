package runtime

import "fmt"

type String struct {
	Value string
}

func (s String) String() string {
	return s.Value
}

type Number struct {
	Value float64
}

func (n Number) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type Boolean struct {
	Value bool
}

func (b Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}
