package parse

import (
	"encoding/json"
	"errors"
)

type Node interface {
	json.Marshaler
	Line() int
	Column() int
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

// Embedded struct to record a node's position
// in the source.
type position struct {
	line   int
	column int
}

func (p position) Line() int {
	return p.line
}

func (p position) Column() int {
	return p.column
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

func (root RootNode) Line() int {
	return 1
}

func (root RootNode) Column() int {
	return 1
}

func (root RootNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(root.Statements)
}

func (root *RootNode) PushStatement(statement *Statement) {
	root.Statements = append(root.Statements, statement)
}

type Statement struct {
	ParentNode
	position
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
	position
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
	position
}

func (i Identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Identifier)
}

type String struct {
	Value string
	position
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
	position
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
	position
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
	position
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
	position
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
		Type:       "declaration",
		ValueType:  *let.Type,
		Identifier: *let.Identifier,
		Children:   let.children,
	})
}

type Assignment struct {
	ParentNode
	Identifier *Identifier
	position
}

func (assignment Assignment) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string
		Identifier Identifier
		Children   []Node
	}{
		Type:       "assignment",
		Identifier: *assignment.Identifier,
		Children:   assignment.children,
	})
}

func (assignment *Assignment) push(child Node) (error, bool) {
	if assignment.Identifier == nil {
		if identifier, isIdentifier := child.(*Identifier); !isIdentifier {
			return errors.New("First child of an assignment statement must be an identifier"), false
		} else {
			assignment.Identifier = identifier
			return nil, true
		}
	}

	return assignment.ParentNode.push(child)
}

type Group struct {
	ParentNode
	position
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

type If struct {
	condition Node
	position
	ParentNode
}

func (i If) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string
		Condition Node
		Children  []Node
	}{
		Type:      "if",
		Condition: i.condition,
		Children:  i.children,
	})
}

func (i *If) push(child Node) (error, bool) {
	if nil == i.condition {
		i.condition = child
	} else {
		return i.ParentNode.push(child)
	}

	return nil, true
}

func (i If) Condition() Node {
	return i.condition
}

func NewStatement(line int, column int, children ...Node) *Statement {
	return &Statement{ParentNode: ParentNode{children: children}, position: position{line: line, column: column}}
}

func NewFunctionCall(identifier string, line int, column int, children ...Node) *FunctionCall {
	return &FunctionCall{Identifier: NewIdentifier(identifier, line, column), ParentNode: ParentNode{children: children}}
}

func NewIdentifier(identifier string, line int, column int) *Identifier {
	return &Identifier{Identifier: identifier, position: position{line: line, column: column}}
}

func NewOperator(operator string, line int, column int, children ...Node) *Operator {
	return &Operator{Operator: operator, ParentNode: ParentNode{children: children}, position: position{line: line, column: column}}
}

func NewString(value string, line int, column int) *String {
	return &String{Value: value, position: position{line: line, column: column}}
}

func NewBoolean(value bool, line int, column int) *Boolean {
	return &Boolean{Value: value, position: position{line: line, column: column}}
}

func NewNumber(value float64, line int, column int) *Number {
	return &Number{Value: value, position: position{line: line, column: column}}
}

func NewDeclaration(identifier Identifier, typeIdentifier Identifier, line int, column int, children ...Node) *Let {
	return &Let{
		children:   children,
		Identifier: &identifier,
		Type:       &typeIdentifier,
		position:   position{line: line, column: column},
	}
}

func NewGroup(line int, column int) *Group {
	return &Group{position: position{line: line, column: column}}
}

func NewAssignment(line int, column int) *Assignment {
	return &Assignment{position: position{line: line, column: column}}
}

func NewIf(line int, column int) *If {
	return &If{position: position{line: line, column: column}}
}
