package dfa

type MachineBuilder interface {
	Path(string, string, string)
	WhenEntering(string, func()) error
	Accept(string) error
	Start(string) (Machine, error)
}

type machineBuilder struct {
	machine *machine
}

func NewMachineBuilder() MachineBuilder {
	return &machineBuilder{newMachine()}
}

func (builder *machineBuilder) Path(from string, via string, to string) {
	if _, exists := builder.machine.states[from]; !exists {
		builder.machine.states[from] = newState(from)
	}

	if _, exists := builder.machine.states[to]; !exists {
		builder.machine.states[to] = newState(to)
	}

	builder.machine.states[from].paths[via] = builder.machine.states[to]
}

func (builder *machineBuilder) Accept(what string) error {
	if err := validateState(builder.machine, what); err != nil {
		return err
	}

	builder.machine.states[what].acceptable = true

	return nil
}

func (builder *machineBuilder) WhenEntering(where string, do func()) error {
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

	return builder.machine, nil
}
