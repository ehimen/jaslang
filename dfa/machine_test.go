package dfa_test

import (
	"testing"

	"github.com/ehimen/jaslang/dfa"
)

func TestAddState(t *testing.T) {
	machine := getMachineBuilder()

	machine.Path("from", "to")
}

func TestWhenEnteringFailsIfStateNotExists(t *testing.T) {
	err := getMachineBuilder().WhenEntering("not-exists", func() {})

	if err == nil {
		t.Error("expected WhenEntering() to fail on non-existent state, but it did not")
	}
}

func TestAcceptFailsIfStateNotExists(t *testing.T) {
	err := getMachineBuilder().Accept("not-exists")

	if err == nil {
		t.Error("expected Accept() to fail on non-existent state, but it did not")
	}
}

func TestAcceptTwoState(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("one", "two")
	builder.Accept("two")

	machine := build(builder, "one", t)

	machine.Transition("two")

	if err := machine.Finish(); err != nil {
		t.Errorf("Expected machine acceptance, but it failed: %v", err)
	}
}

func TestFailTwoState(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("one", "two")

	machine := build(builder, "one", t)

	machine.Transition("two")

	err := machine.Finish()
	if _, expected := err.(dfa.UnacceptableMachineFinishState); !expected {
		t.Error("Expected machine to fail, but it accepted")
	}
}

func TestStartFailsIfStateNotExists(t *testing.T) {
	_, err := getMachineBuilder().Start("not-exists")

	if err == nil {
		t.Error("expected Start() to fail on non-existent state, but it did not")
	}
}

func TestInvalidTransition(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("one", "two")

	machine := build(builder, "one", t)

	if err := machine.Transition("not-exists"); err == nil {
		t.Error("Expected Transition() to fail on non-existent path")
	}
}

func TestWhenEnteringIsCalled(t *testing.T) {
	n := 0

	inc := func() {
		n++
	}

	builder := getMachineBuilder()
	builder.Path("from", "to")
	builder.WhenEntering("to", inc)
	builder.WhenEntering("to", inc)
	builder.WhenEntering("to", inc)

	machine := build(builder, "from", t)

	machine.Transition("to")

	if n != 3 {
		t.Errorf("Expected WhenEntering() to invoke functions, but it didn't. n: %d", n)
	}
}

func TestMachineTrace(t *testing.T) {
	trace := ""

	getTraceFn := func(str string) func() {
		return func() {
			trace += " " + str
		}
	}

	builder := getMachineBuilder()

	builder.Path("one", "two")
	builder.Path("two", "three")
	builder.Path("three", "five")
	builder.Path("four", "five")
	builder.Path("five", "four")
	builder.Path("five", "one")
	builder.WhenEntering("one", getTraceFn("one"))
	builder.WhenEntering("two", getTraceFn("two"))
	builder.WhenEntering("three", getTraceFn("three"))
	builder.WhenEntering("four", getTraceFn("four"))
	builder.WhenEntering("five", getTraceFn("five"))
	builder.Accept("one")

	expected := " two three five four five one"

	machine := build(builder, "one", t)

	machine.Transition("two")
	machine.Transition("three")
	machine.Transition("five")
	machine.Transition("four")
	machine.Transition("five")
	machine.Transition("one")

	if err := machine.Finish(); err != nil {
		t.Errorf("Expected machine to accept, but it didn't: %v", err)
	}

	if expected != trace {
		t.Errorf("Expected trace '%s' but got '%s'", expected, trace)
	}
}

func TestMachineCannotBeUsedAfterCompletion(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("origin", "origin")

	builder.Accept("origin")

	machine := build(builder, "origin", t)

	if err := machine.Finish(); err != nil {
		t.Errorf("Expected machine to accept, but it didn't: %v", err)
	}

	if err := machine.Transition("origin"); err == nil {
		t.Error("Expected second Transition() to fail as machine unusable, but it didn't")
	}

	if err := machine.Finish(); err == nil {
		t.Error("Expected second Finish() to fail as machine unusable, but it didn't")
	}
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
