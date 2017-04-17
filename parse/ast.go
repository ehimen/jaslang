package parse

type Node interface {
}

type ContainsChildren interface {
	Push(child Node)
}

// Can be embedded in to all node types that
// have children.
type ParentNode struct {
	Children []Node
}

func (parent *ParentNode) Push(child Node) {
	parent.Children = append(parent.Children, child)
}

type RootNode struct {
	Statements []*Statement
}

func (root *RootNode) PushStatement(statement *Statement) {
	root.Statements = append(root.Statements, statement)
}

type Statement struct {
	ParentNode
}

type FunctionCall struct {
	Identifier string
	ParentNode
}

type StringLiteral struct {
	Value string
}
