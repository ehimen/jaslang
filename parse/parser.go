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
	nodeStack      []ContainsChildren
	statementStack []Statement
	operators      *Register
	ast            *RootNode
	openedFunction bool
}

type UnexpectedTokenError struct {
	Lexeme lex.Lexeme
	Debug  string
}

func (err UnexpectedTokenError) Error() string {
	return fmt.Sprintf(
		"Unexpected token \"%s\" at position %d, line %d",
		err.Lexeme.Value,
		err.Lexeme.Start,
		err.Lexeme.Line,
	)
}

type InvalidNumberError struct {
	UnexpectedTokenError
}

func (err InvalidNumberError) Error() string {
	return fmt.Sprintf("Invalid number token \"%s\" at position %d", err.Lexeme.Value, err.Lexeme.Start)
}

var UnterminatedStatement = errors.New("Unterminated statement!")

func NewParser(lexer lex.Lexer) Parser {
	parser := parser{lexer: lexer, operators: NewRegister(), openedFunction: false}

	parser.operators.Register("+", 0)
	parser.operators.Register("-", 0)
	parser.operators.Register("*", 1)
	parser.operators.Register("/", 1)
	parser.operators.Register("&&", 0)
	parser.operators.Register("||", 0)
	parser.operators.Register(">", 1)
	parser.operators.Register("<", 1)
	parser.operators.Register("==", 1)

	machine, err := buildDfa(&parser)

	if err != nil {
		panic(fmt.Sprintf("Cannot build parse machine: %v", err))
	}

	parser.dfa = machine

	return &parser
}

func (p *parser) Parse() (RootNode, error) {
	root := &RootNode{}
	p.ast = root
	p.nodeStack = []ContainsChildren{}

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
				if _, isInvalid := err.(dfa.InvalidMachineTransition); isInvalid {
					return *root, UnexpectedTokenError{Lexeme: p.current, Debug: p.dfa.DebugRoute()}
				}

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
		p.openedFunction = true
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
		return InvalidNumberError{UnexpectedTokenError{Lexeme: p.current}}
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
	if len(p.nodeStack) > 0 {
		p.nodeStack = p.nodeStack[0 : len(p.nodeStack)-1]
	}

	return nil
}

// Closes all nodes up the stack, until a group or
// function call node are reached, which we close
// before returning.
func (p *parser) closeGroupOrFunction() error {
	for len(p.nodeStack) > 0 {
		context := getContext(p)

		_, isGroup := context.(*Group)
		_, isFunctionCall := context.(*FunctionCall)

		if isGroup || isFunctionCall {
			p.closeNode()
			break
		} else {
			p.closeNode()
		}
	}

	return nil
}

// Closes a node in an argument list. This will
// close all nodes until it reaches a function call.
// If we're not in a function call, this will return
// an unexpected token error.
func (p *parser) closeArgument() error {
	containsFunctionCall := p.nodeStackContains(func(node ContainsChildren) bool {
		_, isFunctionCall := node.(*FunctionCall)

		return isFunctionCall
	})

	if !containsFunctionCall {
		return UnexpectedTokenError{Lexeme: p.current}
	}

	for _, isFunctionCall := getContext(p).(*FunctionCall); len(p.nodeStack) > 0 && !isFunctionCall; {
		p.closeNode()
	}

	return nil
}

// Creates a precedence grouping. I.e.
// operations enclosed in parentheses.
// Note this will do nothing if we have just
// opened a function call (as this isn't a group).
func (p *parser) createGroup() error {
	if p.openedFunction {
		p.openedFunction = false

		return nil
	}

	p.push(NewGroup())

	return nil
}

// Creates an assignment node for when we have
// assignment without declaration (outside of
// a let).
//
// Note: this won't do anything if we're inside
// a let to avoid us creating this twice.
//
// TODO: better to model declaration and assignment
// TODO: separately?
func (p *parser) createAssignment() error {
	inAssignment := p.nodeStackContains(func(node ContainsChildren) bool {
		_, isAssignment := node.(*Let)

		return isAssignment
	})

	if inAssignment {
		return nil
	}

	p.push(&Assignment{})

	return nil
}

// Closes all open nodes up the stack
// until we close a statement node.
func (p *parser) closeStatement() error {
	context := getContext(p)

	for _, isStatement := context.(*Statement); !isStatement && len(p.nodeStack) > 1; {
		// Close while the head of the stack isn't a statement.
		//println("Closing because not statement")
		p.closeNode()
	}

	// Close one more time. This closes the statement
	p.closeNode()

	return nil
}

func (p *parser) push(node Node) error {
	context := getContext(p)

	// Insert a statement if we need to.
	if context == nil {
		statement := &Statement{}
		p.ast.PushStatement(statement)
		p.nodeStack = append(p.nodeStack, statement)
	}

	if nodeContainingChildren, nodeContainsChildren := node.(ContainsChildren); nodeContainsChildren {
		// Loop over context up the AST until we:
		// 1. Find a context we should should replace.
		// 2. Run out of adjustable AST; simply put it as a child
		//    of the context before we started this loop.
		nodeStackPosition := len(p.nodeStack) - 1

		for {

			if nodeStackPosition < 0 {
				// Run out of AST
				break
			}

			toReplace := p.nodeStack[nodeStackPosition]

			if adjustableParent, isAdjustable := toReplace.(Adjustable); isAdjustable {

				if priority := p.shouldReplaceLastChildOf(nodeContainingChildren, adjustableParent); priority != nil {
					lastChild := adjustableParent.getLastChild()

					// Take the last child of our parent
					adjustableParent.removeLastChild()

					// Push this child on to our new node.
					if err, _ := priority.push(lastChild); err != nil {
						return err
					}

					// Strip the node stack back to the current parent.
					// Those stripped have become children of our
					// current node.
					p.nodeStack = p.nodeStack[0 : nodeStackPosition+1]

					// Now all that remains is to push our new node to
					// the new parent and add our new node to the
					// stack. This all happens as normal outside of
					// this loop.
					break

				} else {
					nodeStackPosition--
				}
			} else {
				// Hit a non-adjustable parent in the AST; stop trying to replace.
				break
			}
		}
	}

	context = getContext(p)

	if err, _ := context.push(node); err != nil {
		return err
	}

	if parent, isParent := node.(ContainsChildren); isParent {
		p.nodeStack = append(p.nodeStack, parent)
	}

	return nil
}

// Returns replacer if it should replace the last child
// of parent.
// This is true when replacer and parent are operators,
// and replacer has a higher precedence. This ensures that
// replacer appears lower in the AST, thus it is evaluated
// first (e.g. "1 + 2 * 3" is "1 + (2 * 3), replacer is *,
// parent is +).
//
// Also true when parent's last child and replacer are
// operators, but replacer does not takes precedence
// over the last child. Ensures that operators are chained
// and the first operator lower in the AST, thus evaluated
// first (e.g. "1 + 2 + 3" is "(1 + 2) + 3", replacer is
// second +, parent is statement whose last child is first +).
//
// Also true when parent's last child does not contain
// children and our replacer is an operator. Ensures that lone
// nodes are placed under an operator as an operator (for now)
// always takes a LHS argument.
//
// Also true when parent's last child is a group and the replacer
// is an operator. This ensures that the group becomes a child
// of the operator, and that the forced precedence of the group
// is respected. E.g. (4 + 4) / 2 should have the / operator
// at the top, with two children: the group 4 + 4 and number 2.
//
// Also true when the parent's last child an identifier and
// the replacer is the assignment operation. This means we
// reshuffle the AST so that the assignment becomes the parent
// node, and everything sitting to the right of the "=" become
// children of the assignment.
func (p parser) shouldReplaceLastChildOf(replacer ContainsChildren, parent Adjustable) ContainsChildren {

	// Check operator precedence.
	parentOperator, parentIsOperator := parent.(*Operator)
	replacerOperator, replacerIsOperator := replacer.(*Operator)

	if parentIsOperator && replacerIsOperator {
		takesPrecedence, err := p.operators.TakesPrecedence(replacerOperator.Operator, parentOperator.Operator)
		// TODO: error checking

		// If the parent is an operator and replacer
		// takes precedence over it, replacer should
		// take the last child of parent.
		if takesPrecedence && err == nil {
			return replacer
		} else {
			return nil
		}
	}

	lastChild := parent.getLastChild()
	lastChildOperator, lastChildIsOperator := lastChild.(*Operator)

	if lastChildIsOperator && replacerIsOperator {
		takesPrecedence, err := p.operators.TakesPrecedence(replacerOperator.Operator, lastChildOperator.Operator)
		// TODO: error checking

		// If our replacer does not take precedence over
		// parent's last child, it should replace the parent.
		if !takesPrecedence && err == nil {
			return replacer
		} else {
			return nil
		}
	}

	_, replacerIsAssignment := replacer.(*Assignment)

	if _, lastChildIsIdentifier := lastChild.(*Identifier); replacerIsAssignment && lastChildIsIdentifier {
		return replacer
	}

	// All other cases require our replacer to be an operator.
	if !replacerIsOperator {
		return nil
	}

	if _, lastChildIsGroup := lastChild.(*Group); lastChildIsGroup {
		return replacer
	}

	// If our parent's last child is not an operator
	// and does not contain children, and our replacer is an
	// operator, we should replace.
	if _, lastChildContainsChildren := lastChild.(ContainsChildren); !lastChildContainsChildren {
		return replacer
	}

	return nil
}

// Gets the current head of the node stack.
func getContext(p *parser) ContainsChildren {
	if len(p.nodeStack) == 0 {
		return nil
	}

	return p.nodeStack[len(p.nodeStack)-1]
}

// Returns true if any node exists in the parser's current
// stack where the match function returns true
func (p parser) nodeStackContains(match func(ContainsChildren) bool) bool {
	for _, node := range p.nodeStack {
		if match(node) {
			return true
		}
	}

	return false
}
