package dfa

import (
	"errors"
	"fmt"
)

type state struct {
	name         string
	whenEntering []func() error
	paths        map[string]*state
	acceptable   bool
}

type Machine interface {
	Transition(string) error
	Finish() error
}

type machine struct {
	current  *state
	states   map[string]*state
	finished bool
}

var MachineUnusable = errors.New("Machine is in an unusable state. Has it already finished?")

type UnacceptableMachineFinishState error

type UnknownMachineState error

type InvalidMachineTransition error

func newMachine() *machine {
	return &machine{states: make(map[string]*state), finished: false}
}

func newState(name string) *state {
	return &state{
		name,
		make([]func() error, 0),
		make(map[string]*state),
		false,
	}
}

func (machine *machine) Transition(how string) error {
	if machine.finished {
		return MachineUnusable
	}

	if _, exists := machine.current.paths[how]; !exists {
		return errors.New(fmt.Sprintf("Don't know how to move from %s to %s", machine.current.name, how))
	}

	machine.current = machine.current.paths[how]

	// Call all functions as we enter this new state
	for _, fn := range machine.current.whenEntering {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (machine *machine) Finish() error {
	if machine.finished {
		return MachineUnusable
	}

	if machine.current.acceptable {
		machine.finished = true
		return nil
	}

	return UnacceptableMachineFinishState(errors.New(fmt.Sprintf("Unacceptable finish state: %s", machine.current.name)))
}

func validateState(machine *machine, state string) error {
	if _, exists := machine.states[state]; !exists {
		return UnknownMachineState(errors.New(fmt.Sprintf("Unknown machine state: %s", state)))
	}

	return nil
}
