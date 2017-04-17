package dfa_test

import (
	"testing"

	"github.com/ehimen/jaslang/dfa"
)

func TestAddState(t *testing.T) {
	machine := getMachineBuilder()

	machine.Path("from", "via", "to")
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

	builder.Path("one", "via", "two")
	builder.Accept("two")

	machine, err := builder.Start("one")

	if err != nil {
		t.Fatalf("Unexpected error when starting machine: %v", err)
	}

	machine.Transition("via")
	err = machine.Finish()

	if err != nil {
		t.Errorf("Expected machine acceptance, but it failed: %v", err)
	}
}

func TestFailTwoState(t *testing.T) {
	builder := getMachineBuilder()

	builder.Path("one", "via", "two")

	machine, err := builder.Start("one")

	if err != nil {
		t.Fatalf("Unexpected error when starting machine: %v", err)
	}

	machine.Transition("via")
	err = machine.Finish()

	if err == nil {
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

	builder.Path("one", "via", "two")

	machine, err := builder.Start("one")

	if err != nil {
		t.Fatalf("Unexpected error when starting machine: %v", err)
	}

	err = machine.Transition("not-via")

	if err == nil {
		t.Error("Expected Transition() to fail on non-existent path")
	}
}

func getMachineBuilder() dfa.MachineBuilder {
	return dfa.NewMachineBuilder()
}
