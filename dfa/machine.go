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
	DebugRoute() string
}

type trace struct {
	path  string
	state state
}

type machine struct {
	current     *state
	states      map[string]*state
	finished    bool
	route       []trace
	transitions map[string]func() error
}

type UnacceptableMachineFinishState struct {
	state string
}

func (err UnacceptableMachineFinishState) Error() string {
	return fmt.Sprintf("Unacceptable finish state: %s", err.state)
}

var MachineUnusable = errors.New("Machine is in an unusable state. Has it already finished?")

type UnknownMachineState struct {
	state string
}

func (err UnknownMachineState) Error() string {
	return fmt.Sprintf("Unknown machine state: %s", err.state)
}

type InvalidMachineTransition struct {
	from string
	to   string
}

func (err InvalidMachineTransition) Error() string {
	return fmt.Sprintf("Don't know how to move from %s to %s", err.from, err.to)
}

func newMachine() *machine {
	return &machine{states: make(map[string]*state), finished: false, transitions: make(map[string]func() error)}
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
		return InvalidMachineTransition{machine.current.name, how}
	}

	machine.current = machine.current.paths[how]
	machine.route = append(machine.route, trace{path: how, state: *machine.current})

	if fn, exists := machine.transitions[how]; exists {
		if err := fn(); err != nil {
			return err
		}
	}

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

	return UnacceptableMachineFinishState{machine.current.name}
}

func (machine *machine) DebugRoute() string {
	trace := ""

	for i, element := range machine.route {
		if i == 0 {
			// Start state
			trace = "ORIGIN: " + element.state.name
		} else {
			trace = fmt.Sprintf("%s >>%s>> %s", trace, element.path, element.state.name)
		}
	}

	return trace
}

func validateState(machine *machine, state string) error {
	if _, exists := machine.states[state]; !exists {
		return UnknownMachineState{state}
	}

	return nil
}
