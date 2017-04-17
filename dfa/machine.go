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
	Path(string, string, string)
	WhenEntering(string, func()) error
	Accept(string) error
	Start(string) error
	Transition(string) error
	Finish() error
}

type machine struct {
	current *state
	states  map[string]*state
}

func NewMachine() Machine {
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

func (machine *machine) Path(from string, via string, to string) {
	if _, exists := machine.states[from]; !exists {
		machine.states[from] = newState(from)
	}

	if _, exists := machine.states[to]; !exists {
		machine.states[to] = newState(to)
	}

	machine.states[from].paths[via] = machine.states[to]
}

func (machine *machine) WhenEntering(where string, do func()) error {
	if err := validateState(machine, where); err != nil {
		return err
	}

	machine.states[where].whenEntering = append(machine.states[where].whenEntering, do)

	return nil
}

func (machine *machine) Accept(what string) error {
	if err := validateState(machine, what); err != nil {
		return err
	}

	machine.states[what].acceptable = true

	return nil
}

func (machine *machine) Start(where string) error {
	if err := validateState(machine, where); err != nil {
		return err
	}

	machine.current = machine.states[where]

	return nil
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
