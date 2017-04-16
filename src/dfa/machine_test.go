package dfa_test

import (
	"dfa"
	"testing"
)

func TestNewMachine(t *testing.T) {
	machine := getMachine()

	if _, ok := machine.(dfa.Machine); !ok {
		t.Error("not a machine")
	}
}

func TestAddState(t *testing.T) {
	machine := getMachine()

	machine.Path("from", "via", "to")
}

func TestWhenEnteringFailsIfStateNotExists(t *testing.T) {
	err := getMachine().WhenEntering("not-exists", func() {})

	if err == nil {
		t.Error("expected WhenEntering() to fail on non-existent state, but it did not")
	}
}

func TestAcceptFailsIfStateNotExists(t *testing.T) {
	err := getMachine().Accept("not-exists")

	if err == nil {
		t.Error("expected Accept() to fail on non-existent state, but it did not")
	}
}

func TestAcceptTwoState(t *testing.T) {
	machine := getMachine()

	machine.Path("one", "via", "two")
	machine.Accept("two")

	machine.Start("one")
	machine.Transition("via")
	err := machine.Finish()

	if err != nil {
		t.Errorf("Expected machine acceptance, but it failed: %v", err)
	}
}

func TestFailTwoState(t *testing.T) {
	machine := getMachine()

	machine.Path("one", "via", "two")

	machine.Start("one")
	machine.Transition("via")
	err := machine.Finish()

	if err == nil {
		t.Error("Expected machine to fail, but it accepted")
	}
}

func TestStartFailsIfStateNotExists(t *testing.T) {
	err := getMachine().Start("not-exists")

	if err == nil {
		t.Error("expected Start() to fail on non-existent state, but it did not")
	}
}

func TestInvalidTransition(t *testing.T) {
	machine := getMachine()

	machine.Path("one", "via", "two")

	machine.Start("one")
	err := machine.Transition("not-via")

	if err == nil {
		t.Error("Expected Transition() to fail on non-existent path")
	}
}

func getMachine() dfa.Machine {
	return dfa.NewMachine()
}
