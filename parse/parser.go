package parse

import (
	"fmt"

	"strconv"

	"errors"

	"github.com/ehimen/jaslang/dfa"
	"github.com/ehimen/jaslang/lex"
)

type Parser interface {
	Parse() (RootNode, error)
}

type parser struct {
	lexer          lex.Lexer
	dfa            dfa.Machine
	current        lex.Lexeme
	next           lex.Lexeme
	nodeStack      []Node
	statementStack []Statement
}

type UnexpectedTokenError struct {
	l lex.Lexeme
}

func (err UnexpectedTokenError) Error() string {
	return fmt.Sprintf("Unexpected token \"%s\" at position %d", err.l.Value, err.l.Start)
}

type InvalidNumberError struct {
	UnexpectedTokenError
}

func (err InvalidNumberError) Error() string {
	return fmt.Sprintf("Invalid number token \"%s\" at position %d", err.l.Value, err.l.Start)
}

var UnterminatedStatement = errors.New("Unterminated statement!")

func NewParser(lexer lex.Lexer) Parser {
	parser := parser{lexer: lexer}

	machine, err := buildDfa(&parser)

	if err != nil {
		panic(fmt.Sprintf("Cannot build parse machine: %v", err))
	}

	parser.dfa = machine

	return &parser
}

func (p *parser) Parse() (RootNode, error) {
	root := &RootNode{}

	p.nodeStack = []Node{root}

	if next, eof, err := p.consume(); eof != nil {
		return *root, nil
	} else if err != nil {
		return *root, err
	} else {
		p.current = next
	}

	if next, eof, err := p.consume(); err != nil {
		return *root, err
	} else if eof == nil {
		p.next = next
	}

	for {
		if p.current.IsEmpty() {
			break
		}

		// We don't care about whitespace
		if p.current.Type != lex.LWhitespace {
			if err := p.dfa.Transition(p.current.Type.String()); err != nil {
				return *root, err
			}
		}

		p.current = p.next

		if next, _, err := p.consume(); err != nil {
			return *root, err
		} else {
			p.next = next
		}
	}

	if err := p.dfa.Finish(); err != nil {
		return *root, UnterminatedStatement
	}

	return *root, nil
}

func (p *parser) consume() (next lex.Lexeme, eof error, lexErr error) {
	if lexeme, err := p.lexer.GetNext(); err == nil {
		next = lexeme
	} else if err == lex.EndOfInput {
		eof = err
	} else {
		lexErr = err
	}

	return
}

func (p *parser) createIdentifier() error {
	if p.next.Type == lex.LParenOpen {
		return p.push(NewFunctionCall(p.current.Value))
	} else {
		return p.push(NewIdentifier(p.current.Value))
	}
}

func (p *parser) createStringLiteral() error {
	return p.push(NewString(p.current.Value))
}

func (p *parser) createBooleanLiteral() error {
	return p.push(NewBoolean(p.current.Type == lex.LBoolTrue))
}

func (p *parser) createNumberLiteral() error {
	if number, err := strconv.ParseFloat(p.current.Value, 64); err == nil {
		p.push(NewNumber(number))
	} else {
		return InvalidNumberError{UnexpectedTokenError{p.current}}
	}

	return nil
}

func (p *parser) createOperator() error {
	return p.push(NewOperator(p.current.Value))
}

func (p *parser) createLet() error {
	return p.push(&Let{})
}

func (p *parser) closeNode() error {
	p.nodeStack = p.nodeStack[0 : len(p.nodeStack)-1]

	return nil
}

func (p *parser) push(node Node) error {
	context := getContext(p)

	// Insert a statement if we need to.
	if root, isRoot := context.(*RootNode); isRoot {
		statement := &Statement{}
		root.PushStatement(statement)
		p.nodeStack = append(p.nodeStack, statement)
	}

	context = getContext(p)

	if parent, isParent := context.(ContainsChildren); isParent {

		if adjustable, isAdjustable := parent.(Adjustable); isAdjustable {
			lastChild := adjustable.getLastChild()

			if priority := takesPrecedence(node, lastChild); priority != nil {
				adjustable.removeLastChild()
				if err, _ := priority.push(lastChild); err != nil {
					return err
				}
			}
		}

		if err, _ := parent.push(node); err != nil {
			return err
		}

	}

	if parent, isParent := node.(ContainsChildren); isParent {
		p.nodeStack = append(p.nodeStack, parent)
	}

	return nil
}

// Returns what if over is what is a parent and over
// is not.
// TODO: Might want to introduce operator priorities.
func takesPrecedence(what Node, over Node) ContainsChildren {
	if over == nil {
		return nil
	}

	_, overIsParent := over.(ContainsChildren)

	if what, whatIsParent := what.(ContainsChildren); whatIsParent && !overIsParent {
		return what
	}

	return nil
}

func getContext(p *parser) Node {
	return p.nodeStack[len(p.nodeStack)-1]
}
