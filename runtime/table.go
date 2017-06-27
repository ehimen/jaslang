package runtime

import (
	"fmt"
	"strings"
)

type entry struct {
	identifier string
	value      Value
}

type operatorEntry struct {
	operator  string
	operands  Types
	operation Invokable
}

type SymbolTable struct {
	parent    *SymbolTable
	entries   map[string]*entry
	operators []operatorEntry
}

type UnknownIdentifier struct {
	identifier string
}

func (err UnknownIdentifier) Error() string {
	return fmt.Sprintf("Unknown identifier: %s", err.identifier)
}

type UnknownOperator struct {
	operator string
	operands []Type
}

func (err UnknownOperator) Error() string {
	operandDescription := []string{}

	for _, operand := range err.operands {
		operandDescription = append(operandDescription, string(operand))
	}

	return fmt.Sprintf("Unknown operator %s with operands (%s)", err.operator, strings.Join(operandDescription, ", "))
}

func NewTable() *SymbolTable {
	return &SymbolTable{entries: make(map[string]*entry)}
}

func (table *SymbolTable) AddFunction(identifier string, invokable Invokable) {
	table.entries[identifier] = &entry{identifier: identifier, value: invokable}
}

func (table *SymbolTable) AddOperator(operator string, operands Types, invokable Invokable) {
	table.operators = append(table.operators, operatorEntry{operator: operator, operands: operands, operation: invokable})
}

func (table *SymbolTable) Operator(operator string, operands Types) (Invokable, error) {
	for _, candidate := range table.operators {
		if candidate.operator == operator && operands.Equal(candidate.operands) {
			return candidate.operation, nil
		}
	}

	return nil, UnknownOperator{operator: operator, operands: operands}
}

func (table *SymbolTable) Invokable(identifier string) (Invokable, error) {
	if entry, exists := table.entries[identifier]; exists {
		if invokable, isInvokable := entry.value.(Invokable); isInvokable {
			return invokable, nil
		}
	}

	return nil, UnknownIdentifier{identifier: identifier}
}
