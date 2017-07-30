package dfa

import (
	"errors"
	"fmt"
)

type MachineBuilder interface {
	Path(from string, via string, to string) error
	Paths(froms []string, via string, tos []string) error
	WhenEntering(string, func() error) error
	Accept(state string) error
	Start(state string) (Machine, error)
	WhenTransitioningVia(string, func() error)
}

type machineBuilder struct {
	machine *machine
}

func NewMachineBuilder() MachineBuilder {
	return &machineBuilder{newMachine()}
}

func (builder *machineBuilder) Path(from string, how string, to string) error {
	if _, exists := builder.machine.states[from]; !exists {
		builder.machine.states[from] = newState(from)
	}

	if _, exists := builder.machine.states[to]; !exists {
		builder.machine.states[to] = newState(to)
	}

	if _, exists := builder.machine.states[from].paths[how]; exists {
		// TODO: not panic?!
		//panic(fmt.Sprintf(`Path "%s" already exists from "%s"`, how, from))
		return errors.New(fmt.Sprintf(`Path "%s" already exists from "%s"`, how, from))
	}

	builder.machine.states[from].paths[how] = builder.machine.states[to]

	return nil
}

func (builder *machineBuilder) WhenTransitioningVia(how string, what func() error) {
	builder.machine.transitions[how] = what
}

func (builder *machineBuilder) Paths(from []string, how string, to []string) error {
	for _, f := range from {
		for _, t := range to {
			if err := builder.Path(f, how, t); err != nil {
				return err
			}
		}
	}

	return nil
}

func (builder *machineBuilder) Accept(what string) error {
	if err := validateState(builder.machine, what); err != nil {
		return err
	}

	builder.machine.states[what].acceptable = true

	return nil
}

func (builder *machineBuilder) WhenEntering(where string, do func() error) error {
	if err := validateState(builder.machine, where); err != nil {
		return err
	}

	builder.machine.states[where].whenEntering = append(builder.machine.states[where].whenEntering, do)

	return nil
}

func (builder *machineBuilder) Start(where string) (Machine, error) {
	var machine Machine

	if err := validateState(builder.machine, where); err != nil {
		return machine, err
	}

	builder.machine.current = builder.machine.states[where]
	builder.machine.route = append(builder.machine.route, trace{state: *builder.machine.current})

	return builder.machine, nil
}
