package lexer

import (
	"errors"
	"fmt"
)

type LexemeType string

const (
	LQuoted     LexemeType = "quoted"
	LIdentifier LexemeType = "identifier"
	LWhitespace LexemeType = "whitespace"
)

type Lexeme struct {
	Start int
	Type  LexemeType
	Value string
}

func (l Lexeme) String() string {
	return fmt.Sprintf("\"%s\" (%s) at position %d", l.Value, l.Type, l.Start)
}

type Lexer interface {
	GetNext() (Lexeme, error)
}

var EndOfInput = errors.New("End of Input")