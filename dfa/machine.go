package dfa

import (
	"errors"
	"fmt"
)

type state struct {
	name         string
	whenEntering []func()
	paths        map[string]*state
	acceptable   bool
}

type Machine interface {
	Transition(string) error
	Finish() error
}

type machine struct {
	current *state
	states  map[string]*state
}

func newMachine() *machine {
	return &machine{states: make(map[string]*state)}
}

func newState(name string) *state {
	return &state{
		name,
		make([]func(), 0),
		make(map[string]*state),
		false,
	}
}

func (machine *machine) Transition(how string) error {
	if _, exists := machine.current.paths[how]; !exists {
		return errors.New(fmt.Sprintf("Don't know how to move from %s via path %s", machine.current.name, how))
	}

	machine.current = machine.current.paths[how]

	return nil
}

func (machine *machine) Finish() error {
	if machine.current.acceptable {
		return nil
	}

	return errors.New(fmt.Sprintf("Cannot accept final state %s", machine.current.name))
}

func validateState(machine *machine, state string) error {
	if _, exists := machine.states[state]; !exists {
		return errors.New(fmt.Sprintf("Unknown state %s", state))
	}

	return nil
}
