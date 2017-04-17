package parse

type Node interface {
}

type ParentNode interface {
	Push(node Node)
}

type RootNode struct {
	Statements []*Statement
}

func (root *RootNode) PushStatement(statement *Statement) {
	root.Statements = append(root.Statements, statement)
}

type Statement struct {
	Children []Node
}

func (statement *Statement) Push(node Node) {
	statement.Children = append(statement.Children, node)
}

type FunctionCall struct {
	Identifier string
	Children   []Node
}

func (functionCall *FunctionCall) Push(node Node) {
	functionCall.Children = append(functionCall.Children, node)
}

type StringLiteral struct {
	Value string
}
