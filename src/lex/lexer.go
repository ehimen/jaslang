package lex

import (
	"errors"
	"fmt"
)

type LexemeType string

const (
	LQuoted     LexemeType = "quoted"
	LIdentifier LexemeType = "identifier"
	LWhitespace LexemeType = "whitespace"
	LParenOpen  LexemeType = "paren-open"
	LParenClose LexemeType = "paren-close"
	LBraceOpen  LexemeType = "brace-open"
	LBraceClose LexemeType = "brace-close"
	LSemiColon  LexemeType = "semi-colon"
	LNumber     LexemeType = "number"
	LOperator   LexemeType = "operator"
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
var UnterminatedString = errors.New("Unterminated string")

type UnexpectedToken error

func NewUnexpectedToken(token string) UnexpectedToken {
	return errors.New(fmt.Sprintf("Unexpected token: %s", token))
}
