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
			if err := p.dfa.Transition(string(p.current.Type)); err != nil {
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

func NewParser(lexer lex.Lexer) Parser {
	parser := parser{lexer: lexer}

	builder := dfa.NewMachineBuilder()

	start := "start"
	identifier := string(lex.LIdentifier)
	parenOpen := string(lex.LParenOpen)
	parenClose := string(lex.LParenClose)
	quoted := string(lex.LQuoted)
	term := string(lex.LSemiColon)
	number := string(lex.LNumber)
	true := string(lex.LBoolTrue)
	false := string(lex.LBoolFalse)
	operator := string(lex.LOperator)
	let := string(lex.LLet)

	literals := []string{number, quoted, true, false}

	builder.Paths([]string{start}, append(literals, identifier, let))
	builder.Paths([]string{identifier}, []string{parenOpen, operator, term})
	builder.Path(parenOpen, quoted)
	builder.Path(quoted, parenClose)
	builder.Path(parenClose, term)
	builder.Paths(literals, []string{term, operator})
	builder.Paths([]string{term}, literals)
	builder.Paths([]string{operator}, append(literals, identifier))
	builder.Path(let, identifier)

	builder.WhenEntering(identifier, parser.createIdentifier)
	builder.WhenEntering(quoted, parser.createStringLiteral)
	builder.WhenEntering(parenClose, parser.closeNode)
	builder.WhenEntering(term, parser.closeNode)
	builder.WhenEntering(number, parser.createNumberLiteral)
	builder.WhenEntering(true, parser.createBooleanLiteral)
	builder.WhenEntering(false, parser.createBooleanLiteral)
	builder.WhenEntering(operator, parser.createOperator)
	builder.WhenEntering(let, parser.createLet)

	builder.Accept(term)

	machine, err := builder.Start(start)

	if err != nil {
		panic(fmt.Sprintf("Cannot build parse machine: %v", err))
	}

	parser.dfa = machine

	return &parser
}

func (p *parser) createIdentifier() error {
	if p.next.Type == lex.LParenOpen {
		p.push(NewFunctionCall(p.current.Value))
	} else {
		p.push(NewIdentifier(p.current.Value))
	}

	return nil
}

func (p *parser) createStringLiteral() error {
	p.push(NewString(p.current.Value))

	return nil
}

func (p *parser) createBooleanLiteral() error {
	p.push(NewBoolean(p.current.Type == lex.LBoolTrue))

	return nil
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
	p.push(NewOperator(p.current.Value))

	return nil
}

func (p *parser) createLet() error {
	p.push(&Let{})

	return nil
}

func (p *parser) closeNode() error {
	p.nodeStack = p.nodeStack[0 : len(p.nodeStack)-1]

	return nil
}

func (p *parser) push(node Node) {
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
				priority.push(lastChild)
			}
		}

		parent.push(node)
	}

	if parent, isParent := node.(ContainsChildren); isParent {
		p.nodeStack = append(p.nodeStack, parent)
	}
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
