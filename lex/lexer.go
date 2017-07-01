package lex

import (
	"errors"
	"fmt"
)

type LexemeType string

func (ltype LexemeType) is(str string) bool {
	return str == string(ltype)
}

func (ltype LexemeType) String() string {
	return string(ltype)
}

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
	LIf         LexemeType = "if"
	LElse       LexemeType = "else"
	LElseIf     LexemeType = "elseif"
	LWhile      LexemeType = "while"
	LLet        LexemeType = "let"
	LBoolTrue   LexemeType = "true"
	LBoolFalse  LexemeType = "false"
	LEquals     LexemeType = "="
	LComma      LexemeType = ","

	OperatorSymbols   string = "+-.^*&/|=>!"
	SpecialCharacters string = "{}();,"
)

var Keywords = []LexemeType{LIf, LElse, LElseIf, LLet, LWhile}

type Lexeme struct {
	Start int
	Line  int
	Type  LexemeType
	Value string
}

func (l Lexeme) String() string {
	return fmt.Sprintf("\"%s\" (%s) at position %d:%d", l.Value, l.Type, l.Line, l.Start)
}

func (l Lexeme) IsEmpty() bool {
	return l == Lexeme{}
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
