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

type String struct {
	Value string
}

type Boolean struct {
	Value bool
}

type Number struct {
	Value float64
}

func NewStatement(children ...Node) *Statement {
	return &Statement{ParentNode: ParentNode{Children: children}}
}

func NewFunctionCall(identifier string, children ...Node) *FunctionCall {
	return &FunctionCall{Identifier: identifier, ParentNode: ParentNode{Children: children}}
}

func NewString(value string) *String {
	return &String{value}
}

func NewBoolean(value bool) *Boolean {
	return &Boolean{value}
}

func NewNumber(value float64) *Number {
	return &Number{value}
}
