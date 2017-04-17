package parse

type Node struct {
}

type Parser interface {
	Parse() Node
}
