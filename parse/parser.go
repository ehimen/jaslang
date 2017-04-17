package parse

import (
	"fmt"

	"github.com/ehimen/jaslang/dfa"
	"github.com/ehimen/jaslang/lex"
)

type Parser interface {
	Parse() RootNode
}

type parser struct {
	lexer          lex.Lexer
	dfa            dfa.Machine
	current        lex.Lexeme
	nodeStack      []Node
	statementStack []Statement
}

func (p *parser) Parse() RootNode {
	root := &RootNode{}

	p.nodeStack = []Node{root}

	for {
		lexeme, err := p.lexer.GetNext()

		if err != nil {
			break
		}

		p.current = lexeme

		p.dfa.Transition(string(lexeme.Type))
	}

	return *root
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

	builder.Path(start, identifier)
	builder.Path(identifier, parenOpen)
	builder.Path(parenOpen, quoted)
	builder.Path(quoted, parenClose)
	builder.Path(parenClose, term)

	builder.WhenEntering(identifier, parser.createIdentifier)
	builder.WhenEntering(quoted, parser.createStringLiteral)
	builder.WhenEntering(parenClose, parser.closeNode)

	builder.Accept(term)

	machine, err := builder.Start(start)

	if err != nil {
		panic(fmt.Sprintf("Cannot build parse machine: %v", err))
	}

	parser.dfa = machine

	return &parser
}

func (p *parser) createIdentifier() {
	p.push(&FunctionCall{Identifier: p.current.Value})
}

func (p *parser) createStringLiteral() {
	p.push(&StringLiteral{Value: p.current.Value})
}

func (p *parser) closeNode() {
	p.nodeStack = p.nodeStack[0 : len(p.nodeStack)-1]
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
