package lexer

import (
	"io"
	"strings"
)

type jslLexer struct {
	reader io.RuneScanner
	ch chan Lexeme
	position int
	start int
	fn stateFunction
}

func NewJslLexer(reader io.RuneScanner) Lexer {
	return &jslLexer{reader, make(chan Lexeme, 1), 1, 1, defaultState}
}

func (l *jslLexer) GetNext() (Lexeme, error) {

	var lexeme Lexeme

	for {
		select {
		case next := <- l.ch:
			return next, nil
		default:
			if l.fn == nil {
				return lexeme, EndOfInput
			}

			l.start = l.position

			l.fn = l.fn(l)
		}
	}
}

func (l *jslLexer) next() (string, error) {
	l.position += 1
	next, _, err := l.reader.ReadRune()
	return string(next), err
}

func (l *jslLexer) back() {
	l.position -= 1
	l.reader.UnreadRune()
}

func (l *jslLexer) peek() (next string, err error) {
	next, err = l.next()
	l.back()
	return
}

func (l *jslLexer) emit(value string, lexemeType LexemeType) {
	if len(value) == 0 {
		// Prevents emitting empty tokens
		return
	}

	l.ch <- Lexeme{
		l.start,
		lexemeType,
		value,
	}
}

type stateFunction func(l *jslLexer) stateFunction

func quotedState(l *jslLexer) stateFunction {
	lType := LQuoted
	value := ""

	// Have to move start back to cover preceding quote
	// TODO: can do better?
	l.start--

	for {
		done, current := isDone(l, value, lType)

		if done {
			return nil
		}

		if isQuote(current) {
			l.emit(value, lType)
			l.back()
			return defaultState
		}

		value += current
	}
}

func defaultState(l *jslLexer) stateFunction {
	next, err := l.peek()

	if err != nil {
		return nil
	}

	switch {
	case isQuote(next):
		l.next()
		return quotedState
	case isSpace(next):
		return spaceState
	}

	return identifierState
}

func isQuote(str string) bool {
	return str == "\""
}

func isSpace(str string) bool {
	return strings.TrimSpace(str) == ""
}

func spaceState(l *jslLexer) stateFunction {
	lType := LWhitespace

	value := ""

	for {
		done, current := isDone(l, value, lType)

		if done {
			return nil
		}

		if !isSpace(current) {
			l.emit(value, lType)
			l.back()
			return defaultState
		}

		value += current
	}
}

func identifierState(l *jslLexer) stateFunction {
	lType := LIdentifier

	value := ""

	for {
		done, current := isDone(l, value, lType)

		if done {
			return nil
		}

		if isSpace(current) || isQuote(current) {
			l.emit(value, lType)
			l.back()
			return defaultState
		}

		value += current
	}
}

func isDone(l *jslLexer, value string, lexemeType LexemeType) (done bool, next string) {
	cur, err := l.next()

	if err != nil {
		l.emit(value, lexemeType)
		return true, cur
	}

	return false, cur
}