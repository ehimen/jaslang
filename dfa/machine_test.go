package dfa_test

import (
	"testing"

	"errors"

	"github.com/ehimen/jaslang/dfa"
	"github.com/stretchr/testify/assert"
)

func TestAddState(t *testing.T) {
	machine := getMachineBuilder()

	machine.Path("from", "only", "to")
}

func TestWhenEnteringFailsIfStateNotExists(t *testing.T) {
	err := getMachineBuilder().WhenEntering("not-exists", func() error { return nil })

	if _, ok := err.(dfa.UnknownMachineState); !ok {
		t.Error("expected WhenEntering() to fail on non-existent state, but it did not")
	}
}

func TestAcceptFailsIfStateNotExists(t *testing.T) {
	err := getMachineBuilder().Accept("not-exists")

	if _, ok := err.(dfa.UnknownMachineState); !ok {
		t.Error("expected Accept() to fail on non-existent state, but it did not")
	}
}

func TestAcceptTwoState(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("one", "only", "two")
	builder.Accept("two")

	machine := build(builder, "one", t)

	machine.Transition("only")

	if err := machine.Finish(); err != nil {
		t.Errorf("Expected machine acceptance, but it failed: %v", err)
	}
}

func TestFailTwoState(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("one", "only", "two")

	machine := build(builder, "one", t)

	machine.Transition("only")

	if _, expected := machine.Finish().(dfa.UnacceptableMachineFinishState); !expected {
		t.Error("Expected machine to fail, but it accepted or returned the wrong error")
	}
}

func TestStartFailsIfStateNotExists(t *testing.T) {
	_, err := getMachineBuilder().Start("not-exists")

	if _, ok := err.(dfa.UnknownMachineState); !ok {
		t.Error("expected Start() to fail on non-existent state, but it did not")
	}
}

func TestInvalidTransition(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("one", "only", "two")

	machine := build(builder, "one", t)

	if _, ok := machine.Transition("not-exists").(dfa.InvalidMachineTransition); !ok {
		t.Error("Expected Transition() to fail on non-existent path")
	}
}

func TestWhenEnteringIsCalled(t *testing.T) {
	n := 0

	inc := func() error {
		n++

		return nil
	}

	builder := getMachineBuilder()
	builder.Path("from", "only", "to")
	builder.WhenEntering("to", inc)
	builder.WhenEntering("to", inc)
	builder.WhenEntering("to", inc)

	machine := build(builder, "from", t)

	machine.Transition("only")

	if n != 3 {
		t.Errorf("Expected WhenEntering() to invoke functions, but it didn't. n: %d", n)
	}
}

func TestMachineTrace(t *testing.T) {
	trace := ""

	getTraceFn := func(str string) func() error {
		return func() error {
			trace += " " + str

			return nil
		}
	}

	builder := getMachineBuilder()

	builder.Path("one", "one-two", "two")
	builder.Path("two", "two-three", "three")
	builder.Path("three", "three-five", "five")
	builder.Path("four", "four-five", "five")
	builder.Path("five", "five-four", "four")
	builder.Path("five", "five-one", "one")
	builder.WhenEntering("one", getTraceFn("one"))
	builder.WhenEntering("two", getTraceFn("two"))
	builder.WhenEntering("three", getTraceFn("three"))
	builder.WhenEntering("four", getTraceFn("four"))
	builder.WhenEntering("five", getTraceFn("five"))
	builder.Accept("one")

	expected := " two three five four five one"

	machine := build(builder, "one", t)

	machine.Transition("one-two")
	machine.Transition("two-three")
	machine.Transition("three-five")
	machine.Transition("five-four")
	machine.Transition("four-five")
	machine.Transition("five-one")

	if err := machine.Finish(); err != nil {
		t.Errorf("Expected machine to accept, but it didn't: %v", err)
	}

	if expected != trace {
		t.Errorf("Expected trace '%s' but got '%s'", expected, trace)
	}
}

func TestMachineCannotBeUsedAfterCompletion(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("origin", "only", "origin")

	builder.Accept("origin")

	machine := build(builder, "origin", t)

	if err := machine.Finish(); err != nil {
		t.Errorf("Expected machine to accept, but it didn't: %v", err)
	}

	if err := machine.Transition("only"); err == nil {
		t.Error("Expected second Transition() to fail as machine unusable, but it didn't")
	}

	if err := machine.Finish(); err != dfa.MachineUnusable {
		t.Error("Expected second Finish() to fail as machine unusable, but it didn't")
	}
}

func TestWhenFnFailsTransitionFails(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("origin", "only", "origin")

	expected := errors.New("Test error")

	builder.WhenEntering("origin", func() error { return expected })

	machine := build(builder, "origin", t)

	if actual := machine.Transition("only"); actual != expected {
		t.Errorf("Expected Transition() to return same err as callback, but got %v", actual)
	}
}

func TestPaths(t *testing.T) {
	builder := getMachineBuilder()

	builder.Paths([]string{"one", "two"}, "transition", []string{"one", "two"})
	builder.Accept("two")

	machine := build(builder, "one", t)

	machine.Transition("transition")
	machine.Transition("transition")
	machine.Transition("transition")

	if err := machine.Finish(); err != nil {
		t.Errorf("Expected multiple paths to succeed, but it didn't: %v", err)
	}
}

func TestDebugRoute(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("one", "1", "two")
	builder.Path("two", "2", "three")
	builder.Path("three", "3", "one")

	machine := build(builder, "one", t)

	machine.Transition("1")
	machine.Transition("2")
	machine.Transition("3")
	machine.Transition("1")

	assert.Equal(t, "ORIGIN: one >>1>> two >>2>> three >>3>> one >>1>> two", machine.DebugRoute())
}

func build(builder dfa.MachineBuilder, start string, t *testing.T) dfa.Machine {
	machine, err := builder.Start(start)

	if err != nil {
		t.Fatalf("Unexpected error when starting machine: %v", err)
	}

	return machine
}

func getMachineBuilder() dfa.MachineBuilder {
	return dfa.NewMachineBuilder()
}
