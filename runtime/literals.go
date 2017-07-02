package runtime

import "fmt"

type Type string

var TypeNone = Type("none")
var TypeBoolean = Type("boolean")
var TypeNumber = Type("number")
var TypeString = Type("string")
var TypeInvokable = Type("invokable")

func (t Type) DefaultValue() Value {
	switch t {
	case TypeBoolean:
		return Boolean{Value: false}
	case TypeNumber:
		return Number{Value: 0}
	case TypeString:
		return String{Value: ""}
	case TypeInvokable:
		return Noop{}
	}

	return nil
}

type Types []Type

func (types Types) Equal(other Types) bool {
	if len(types) != len(other) {
		return false
	}

	for _, t1 := range types {
		found := false

		for _, t2 := range other {
			if t1 == t2 {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

type String struct {
	Value string
}

func (s String) String() string {
	return s.Value
}

func (s String) Type() Type {
	return TypeString
}

type Number struct {
	Value float64
}

func (n Number) String() string {
	return fmt.Sprintf("%.3f", n.Value)
}

func (n Number) Type() Type {
	return TypeNumber
}

type Boolean struct {
	Value bool
}

func (b Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b Boolean) Type() Type {
	return TypeBoolean
}
