package lexer

import (
	"strings"
)

type jslLexer struct {
	reader *strings.Reader
	ch chan Lexeme
	position int
	start int
	fn stateFunction
}

func NewJslLexer(reader *strings.Reader) Lexer {
	return &jslLexer{reader, make(chan Lexeme, 1), 1, 1, defaultState}
}

func (l *jslLexer) GetNext() (Lexeme, error) {

	for {
		select {
		case next := <- l.ch:
			return next, nil
		default:
			var nothing Lexeme

			if l.fn == nil {
				return nothing, EndOfInput
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

func stringState(l *jslLexer) stateFunction {
	value := ""

	// Move the start position back one to correctly
	// report the start of the string (i.e. the " character)
	// TODO: can we do better?
	l.start--

	for {
		cur, err := l.next()

		if err != nil {
			l.emit(value, LString)
			return nil
		}

		if cur == "\"" {
			l.emit(value, LString)
			return defaultState
		}

		value += cur
	}
}

func defaultState(l *jslLexer) stateFunction {
	value := ""

	for {
		cur, err := l.next()

		if err != nil {
			l.emit(value, LIdentifier)
			return nil
		}

		if cur == "\"" {
			l.emit(value, LIdentifier)
			return stringState
		}

		value += cur
	}
}