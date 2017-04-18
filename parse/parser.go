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
	nodeStack      []Node
	statementStack []Statement
}

// TODO: Better error reporting including failing lexeme
// TODO: Problems with type assertions on error types not
// TODO: working in the way I expected them to.
// TODO: e.g. "type MyError error" doesn't seem to exclusively
// TODO: match when running err.(MyError)
var SyntaxError = errors.New("Syntax error!")

var InvalidNumber = errors.New("Invalid number!")

var UnterminatedStatement = errors.New("Unterminated statement!")

func (p *parser) Parse() (RootNode, error) {
	root := &RootNode{}

	p.nodeStack = []Node{root}

	for {
		lexeme, err := p.lexer.GetNext()

		if err != nil {
			break
		}

		// We don't care about whitespace
		if lexeme.Type == lex.LWhitespace {
			continue
		}

		p.current = lexeme

		if err := p.dfa.Transition(string(lexeme.Type)); err != nil {
			return *root, err
		}
	}

	if err := p.dfa.Finish(); err != nil {
		return *root, UnterminatedStatement
	}

	return *root, nil
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

	// TODO: would be good to group lots of common elements here (e.g. literals)
	builder.Path(start, number)
	builder.Path(start, true)
	builder.Path(start, identifier)
	builder.Path(identifier, parenOpen)
	builder.Path(parenOpen, quoted)
	builder.Path(quoted, parenClose)
	builder.Path(parenClose, term)
	builder.Path(number, term)
	builder.Path(term, number)
	builder.Path(term, false)
	builder.Path(true, term)
	builder.Path(false, term)

	builder.WhenEntering(identifier, parser.createIdentifier)
	builder.WhenEntering(quoted, parser.createStringLiteral)
	builder.WhenEntering(parenClose, parser.closeNode)
	builder.WhenEntering(term, parser.closeNode)
	builder.WhenEntering(number, parser.createNumberLiteral)
	builder.WhenEntering(true, parser.createBooleanLiteral)
	builder.WhenEntering(false, parser.createBooleanLiteral)

	builder.Accept(term)

	machine, err := builder.Start(start)

	if err != nil {
		panic(fmt.Sprintf("Cannot build parse machine: %v", err))
	}

	parser.dfa = machine

	return &parser
}

func (p *parser) createIdentifier() error {
	p.push(&FunctionCall{Identifier: p.current.Value})

	return nil
}

func (p *parser) createStringLiteral() error {
	p.push(&StringLiteral{Value: p.current.Value})

	return nil
}

func (p *parser) createBooleanLiteral() error {
	value := p.current.Type == lex.LBoolTrue

	p.push(&BooleanLiteral{Value: value})

	return nil
}

func (p *parser) createNumberLiteral() error {
	if number, err := strconv.ParseFloat(p.current.Value, 64); err == nil {
		p.push(&NumberLiteral{Value: number})
	} else {
		return InvalidNumber
	}

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
		parent.Push(node)
	}

	if parent, isParent := node.(ContainsChildren); isParent {
		p.nodeStack = append(p.nodeStack, parent)
	}
}

func getContext(p *parser) Node {
	return p.nodeStack[len(p.nodeStack)-1]
}
