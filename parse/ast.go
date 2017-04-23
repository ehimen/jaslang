package parse

type Node interface {
}

type ContainsChildren interface {
	push(child Node)
	getLastChild() Node
	removeLastChild()
}

// Can be embedded in to all node types that
// have children.
type ParentNode struct {
	Children []Node
}

func (parent *ParentNode) push(child Node) {
	parent.Children = append(parent.Children, child)
}

func (parent *ParentNode) getLastChild() Node {
	if len(parent.Children) > 0 {
		return parent.Children[len(parent.Children)-1]
	}

	return nil
}

func (parent *ParentNode) removeLastChild() {
	if len(parent.Children) == 0 {
		return
	}

	parent.Children = parent.Children[0 : len(parent.Children)-1]
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
	Identifier *Identifier
	ParentNode
}

type Identifier struct {
	Identifier string
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

type Operator struct {
	Operator string
	ParentNode
}

func NewStatement(children ...Node) *Statement {
	return &Statement{ParentNode: ParentNode{Children: children}}
}

func NewFunctionCall(identifier string, children ...Node) *FunctionCall {
	return &FunctionCall{Identifier: NewIdentifier(identifier), ParentNode: ParentNode{Children: children}}
}

func NewIdentifier(identifier string) *Identifier {
	return &Identifier{Identifier: identifier}
}

func NewOperator(operator string, children ...Node) *Operator {
	return &Operator{Operator: operator, ParentNode: ParentNode{Children: children}}
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
