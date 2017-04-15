package lex

import (
	"io"
	"strings"
)

type jslLexer struct {
	reader      io.RuneScanner
	ch          chan Lexeme
	position    int
	start       int
	fn          stateFunction
	current     string
	err         error
	lookahead   []string
	initialised bool
	operators   []string
}

func NewJslLexer(reader io.RuneScanner, operators []string) Lexer {
	return &jslLexer{
		reader,
		make(chan Lexeme, 1),
		1,
		1,
		defaultState,
		"",
		nil,
		make([]string, 2),
		false,
		operators,
	}
}

func (l *jslLexer) GetNext() (Lexeme, error) {
	var lexeme Lexeme

	for {
		select {
		case next := <-l.ch:
			return next, nil
		default:
			if l.err != nil {
				return lexeme, l.err
			}

			l.start = l.position

			l.fn, l.err = l.fn(l)
		}
	}
}

func (l *jslLexer) initialise() {
	if l.initialised {
		return
	}

	// Nothing in lookahead. We need to fill it to capacity.
	for i := 0; i < cap(l.lookahead); i++ {
		if nextRune, _, err := l.reader.ReadRune(); nil == err {
			l.lookahead[i] = string(nextRune)
		}
	}

	l.initialised = true
}

func (l *jslLexer) next() string {
	l.initialise()

	next := l.lookahead[0]

	l.current += next

	l.position++

	// Remove the first item and push on the next rune.
	l.lookahead = l.lookahead[1:]
	if nextRune, _, err := l.reader.ReadRune(); nil == err {
		l.lookahead = append(l.lookahead, string(nextRune))
	} else {
		// If we have reached the end of input we still
		// need to fill the lookahead slice to keep its
		// size consistent.
		l.lookahead = append(l.lookahead, "")
	}

	return next
}

func (l *jslLexer) ignore() string {
	next := l.next()
	// Undo the addition of the next rune.
	l.current = l.current[0 : len(l.current)-1]
	// Return the ignored rune
	return next
}

func (l *jslLexer) has() bool {
	return l.peek() != ""
}

func (l *jslLexer) peek() string {
	l.initialise()

	return l.lookahead[0]
}

func (l *jslLexer) peekAhead(amount int) string {
	if amount > cap(l.lookahead) {
		panic("Peek exceeds lookahead")
	}

	return l.lookahead[amount-1]
}

func (l *jslLexer) emit(lexemeType LexemeType) {
	if len(l.current) == 0 {
		return
	}

	l.ch <- Lexeme{
		l.start,
		lexemeType,
		l.current,
	}

	l.current = ""
}

func isQuote(str string) bool {
	return strings.IndexAny(str, "'\"") == 0 && len(str) == 1
}

func isSpace(str string) bool {
	if len(str) == 0 {
		return false
	}

	return strings.TrimSpace(str) == ""
}

func isIdentifierStartCharacter(str string) bool {
	return isIdentifierCharacter(str) && strings.IndexAny(str, "0123456789") != 0
}

func isIdentifierCharacter(str string) bool {
	return !isQuote(str) && !isSpace(str) && !isSpecialSymbol(str)
}

func isSpecialSymbol(next string) bool {
	return strings.IndexAny(next, "{}();") == 0
}

type stateFunction func(l *jslLexer) (stateFunction, error)

func defaultState(l *jslLexer) (stateFunction, error) {
	next := l.peek()

	if next == "" {
		return nil, EndOfInput
	}

	switch {
	case isQuote(next):
		return quotedState, nil
	case isSpace(next):
		return spaceState, nil
	case isSpecialSymbol(next):
		return characterState, nil
	case strings.IndexAny(next, "+-0123456789") == 0:
		return numberState, nil
	case isIdentifierStartCharacter(next):
		return identifierState, nil
	}

	return nil, NewUnexpectedToken(next)
}

func numberState(l *jslLexer) (stateFunction, error) {
	allowSign, doneDecimal := true, false

	for l.has() {
		next := l.peek()

		switch {
		case !doneDecimal && next == ".":
			doneDecimal = true
			l.next()
		case allowSign && strings.IndexAny(next, "+-") == 0:
			nextNext := l.peekAhead(2)

			if strings.IndexAny(nextNext, "0123456789") != 0 {
				// A sign not followed by any digits, this isn't a number
				return identifierState, nil
			}

			l.next()
		case strings.IndexAny(next, "0123456789") == 0:
			l.next()
		default:
			l.emit(LNumber)
			return defaultState, nil
		}

		allowSign = false
	}

	l.emit(LNumber)

	return nil, EndOfInput
}

func characterState(l *jslLexer) (stateFunction, error) {
	next := l.next()

	switch next {
	case "{":
		l.emit(LBraceOpen)
	case "}":
		l.emit(LBraceClose)
	case "(":
		l.emit(LParenOpen)
	case ")":
		l.emit(LParenClose)
	case ";":
		l.emit(LSemiColon)
	}

	return defaultState, nil
}

func quotedState(l *jslLexer) (stateFunction, error) {
	var next, nextNext, quote string
	// Ignore the first quote, and store
	// what type of quote it is
	quote = l.ignore()

	escapableCharacters := "\\" + quote

	for l.has() {
		next = l.peek()

		if next == `\` {
			l.ignore()
			nextNext = l.ignore()

			if strings.IndexAny(nextNext, escapableCharacters) == -1 {
				// Not an escapable character so push both on to the current value.
				l.current += next + nextNext
			} else {
				// Is escapable so ignore the backslash (\) and just put the escaped
				// character on.
				l.current += nextNext
			}

			continue
		}

		if next == quote {
			l.ignore() // Consume the closing quote.
			l.emit(LQuoted)
			return defaultState, nil
		}

		l.next()
	}

	return nil, UnterminatedString
}

func spaceState(l *jslLexer) (stateFunction, error) {
	for l.has() {
		if !isSpace(l.peek()) {
			l.emit(LWhitespace)
			return defaultState, nil
		} else {
			l.next()
		}
	}

	return nil, EndOfInput
}

func identifierState(l *jslLexer) (stateFunction, error) {
	for l.has() {
		next := l.peek()

		if !isIdentifierCharacter(next) {
			l.emit(LIdentifier)
			return defaultState, nil
		}

		l.next()

		for _, operator := range l.operators {
			if operator == l.current {
				l.emit(LOperator)
				return defaultState, nil
			}
		}
	}

	l.emit(LIdentifier)

	return nil, EndOfInput
}
