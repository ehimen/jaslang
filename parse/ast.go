package parse

import (
	"errors"
)

type Node interface {
}

type ContainsChildren interface {
	push(child Node) (error, bool)
	Children() []Node
}

type Adjustable interface {
	ContainsChildren
	getLastChild() Node
	removeLastChild()
}

// Can be embedded in to all node types that
// have children.
type ParentNode struct {
	children []Node
}

func (parent *ParentNode) push(child Node) (error, bool) {
	parent.children = append(parent.children, child)

	return nil, true
}

func (parent *ParentNode) getLastChild() Node {
	if len(parent.children) > 0 {
		return parent.children[len(parent.children)-1]
	}

	return nil
}

func (parent *ParentNode) Children() []Node {
	return parent.children
}

func (parent *ParentNode) removeLastChild() {
	if len(parent.children) == 0 {
		return
	}

	parent.children = parent.children[0 : len(parent.children)-1]
}

// TODO: do we make this ContainsChildren? Would simplify logic,
// TODO: but we lose the strictness of pushing only children.
// TODO: Maybe have a separate AcceptsAnyChild and ContainsChildren
// TODO: interfaces? AcceptsAnyChild.push(), ContainsChildren.Children().
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

type Let struct {
	Identifier *Identifier
	Type       *Identifier
	Children   []Node
}

func (let *Let) push(child Node) (error, bool) {
	if let.Identifier == nil {
		if ident, isIdentifier := child.(*Identifier); !isIdentifier {
			return errors.New("Let requires an identifier"), false // TODO: test this
		} else {
			let.Identifier = ident
			return nil, true
		}
	}

	if let.Type == nil {
		if ident, isIdentifier := child.(*Identifier); !isIdentifier {
			return errors.New("Let requires a type identifier"), false // TODO: test this
		} else {
			let.Type = ident
			return nil, true
		}
	}

	let.Children = append(let.Children, child)

	return nil, true
}

func (let *Let) getLastChild() Node {
	if len(let.Children) > 0 {
		return let.Children[len(let.Children)-1]
	}

	// Deliberately don't return type or identifier; these
	// can't be adjusted.

	return nil
}

func (let *Let) removeLastChild() {
	if len(let.Children) > 0 {
		let.Children = let.Children[0 : len(let.Children)-1]
	}
}

func NewStatement(children ...Node) *Statement {
	return &Statement{ParentNode: ParentNode{children: children}}
}

func NewFunctionCall(identifier string, children ...Node) *FunctionCall {
	return &FunctionCall{Identifier: NewIdentifier(identifier), ParentNode: ParentNode{children: children}}
}

func NewIdentifier(identifier string) *Identifier {
	return &Identifier{Identifier: identifier}
}

func NewOperator(operator string, children ...Node) *Operator {
	return &Operator{Operator: operator, ParentNode: ParentNode{children: children}}
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

func NewLet() *Let {
	return &Let{}
}
