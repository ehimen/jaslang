package runtime

import (
	"fmt"
)

type entry struct {
	identifier string
	//isType      bool
	value Value
}

type SymbolTable struct {
	parent  *SymbolTable
	entries map[string]*entry
}

type UnknownIdentifier struct {
	identifier string
}

func (err UnknownIdentifier) Error() string {
	return fmt.Sprintf("Unknown identifier: %s", err.identifier)
}

func NewTable() *SymbolTable {
	return &SymbolTable{entries: make(map[string]*entry)}
}

func (table *SymbolTable) AddFunction(identifier string, invokable Invokable) {
	table.entries[identifier] = &entry{identifier: identifier, value: invokable}
}

func (table *SymbolTable) Invokable(identifier string) (Invokable, error) {
	if entry, exists := table.entries[identifier]; exists {
		if invokable, isInvokable := entry.value.(Invokable); isInvokable {
			return invokable, nil
		}
	}

	return nil, UnknownIdentifier{identifier: identifier}
}
