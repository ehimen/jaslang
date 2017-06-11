package runtime

type Value interface {
	String() string
	Type() Type
}

type Void struct {
}

func (v Void) String() string {
	return "<void>"
}

func (v Void) Type() Type {
	return TypeNone
}
