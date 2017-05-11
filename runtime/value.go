package runtime

type Value interface {
	String() string
}

type Void struct {
}

func (v Void) String() string {
	return "<void>"
}
