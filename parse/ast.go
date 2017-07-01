package parse

import (
	"encoding/json"
	"errors"
)

type Node interface {
	json.Marshaler
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

func (parent *ParentNode) encodeChildren() ([]byte, error) {
	return json.Marshal(parent.children)
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

func (root RootNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(root.Statements)
}

func (root *RootNode) PushStatement(statement *Statement) {
	root.Statements = append(root.Statements, statement)
}

type Statement struct {
	ParentNode
}

func (s Statement) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string
		Children []Node
	}{
		Type:     "statement",
		Children: s.children,
	})
}

type FunctionCall struct {
	Identifier *Identifier
	ParentNode
}

func (f FunctionCall) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string
		Identifier *Identifier
		Children   []Node
	}{
		Type:       "function",
		Identifier: f.Identifier,
		Children:   f.children,
	})
}

type Identifier struct {
	Identifier string
}

func (i Identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Identifier)
}

type String struct {
	Value string
}

func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string
		Value string
	}{
		Type:  "string",
		Value: s.Value,
	})
}

type Boolean struct {
	Value bool
}

func (b Boolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string
		Value bool
	}{
		Type:  "boolean",
		Value: b.Value,
	})
}

type Number struct {
	Value float64
}

func (n Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string
		Value float64
	}{
		Type:  "number",
		Value: n.Value,
	})
}

type Operator struct {
	Operator string
	ParentNode
}

func (o Operator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Children []Node
		Type     string
		Operator string
	}{
		Children: o.children,
		Type:     "operator",
		Operator: o.Operator,
	})
}

type Let struct {
	Identifier *Identifier
	Type       *Identifier
	children   []Node
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

	let.children = append(let.children, child)

	return nil, true
}

func (let *Let) getLastChild() Node {
	if len(let.children) > 0 {
		return let.children[len(let.children)-1]
	}

	// Deliberately don't return type or identifier; these
	// can't be adjusted.

	return nil
}

func (let *Let) removeLastChild() {
	if len(let.children) > 0 {
		let.children = let.children[0 : len(let.children)-1]
	}
}

func (let *Let) Children() []Node {
	return let.children
}

func (let Let) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string
		ValueType  Identifier
		Identifier Identifier
		Children   []Node
	}{
		Type:       "assignment",
		ValueType:  *let.Type,
		Identifier: *let.Identifier,
		Children:   let.children,
	})
}

//

//type ContainsChildren interface {
//	push(child Node) (error, bool)
//	Children() []Node
//}

type Group struct {
	ParentNode
	//children []Node
}

func (group Group) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string
		Children []Node
	}{
		Type:     "group",
		Children: group.children,
	})
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

func NewLet(identifier string, typeIdentifier string, children ...Node) *Let {
	return &Let{children: children, Identifier: NewIdentifier(identifier), Type: NewIdentifier(typeIdentifier)}
}

func NewGroup() *Group {
	return &Group{}
}
