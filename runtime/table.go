package runtime

import (
	"errors"
	"fmt"
	"strings"
)

type entry struct {
	identifier string
	valueType  Type
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
	types     map[string]Type
}

type UnknownIdentifier struct {
	identifier string
}

func (err UnknownIdentifier) Error() string {
	return fmt.Sprintf("Unknown identifier: %s", err.identifier)
}

type UnknownType struct {
	identifier string
}

func (err UnknownType) Error() string {
	return fmt.Sprintf("Unknown type: %s", err.identifier)
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
	return &SymbolTable{entries: make(map[string]*entry), types: make(map[string]Type)}
}

func (table *SymbolTable) AddType(identifier string, p Type) {
	table.types[identifier] = p
}

func (table *SymbolTable) AddFunction(identifier string, invokable Invokable) {
	table.entries[identifier] = &entry{identifier: identifier, value: invokable, valueType: TypeInvokable}
}

func (table *SymbolTable) AddOperator(operator string, operands Types, invokable Invokable) {
	table.operators = append(table.operators, operatorEntry{operator: operator, operands: operands, operation: invokable})
}

func (table *SymbolTable) Define(identifier string, t Type) error {
	if _, exists := table.entries[identifier]; exists {
		return errors.New(fmt.Sprintf(`Cannot declare symbol "%s"`, identifier))
	}

	// TODO: default value!?

	table.entries[identifier] = &entry{identifier: identifier, valueType: t}

	return nil
}

func (table *SymbolTable) Set(identifier string, value Value) error {

	if valueEntry, exists := table.entries[identifier]; !exists {
		return UnknownIdentifier{identifier: identifier}
	} else {
		if value.Type() != valueEntry.valueType {
			return errors.New(fmt.Sprintf(
				`Invalid value for "%s". Value %s is not of expected type %s`,
				identifier,
				value,
				valueEntry.valueType,
			))
		}

		valueEntry.value = value

		return nil
	}
}

func (table *SymbolTable) Get(identifier string) (Value, error) {
	if entry, exists := table.entries[identifier]; exists {
		return entry.value, nil
	}

	return nil, UnknownIdentifier{identifier: identifier}
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

func (table *SymbolTable) Type(identifier string) (Type, error) {
	var t Type

	if registeredType, exists := table.types[identifier]; exists {
		return registeredType, nil
	}

	return t, UnknownType{identifier: identifier}
}
