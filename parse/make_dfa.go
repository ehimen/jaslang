package parse

import (
	"github.com/ehimen/jaslang/dfa"
	"github.com/ehimen/jaslang/lex"
)

var start = "start"
var identifier = lex.LIdentifier.String()
var parenOpen = lex.LParenOpen.String()
var parenClose = lex.LParenClose.String()
var quoted = lex.LQuoted.String()
var term = lex.LSemiColon.String()
var number = lex.LNumber.String()
var ltrue = lex.LBoolTrue.String()
var lfalse = lex.LBoolFalse.String()
var operator = lex.LOperator.String()
var let = lex.LLet.String()
var equals = lex.LEquals.String()

func buildDfa(p *parser) (dfa.Machine, error) {

	builder := dfa.NewMachineBuilder()

	builder.Path(start, quoted, quoted)
	builder.Path(start, ltrue, ltrue)
	builder.Path(start, lfalse, lfalse)
	builder.Path(start, let, let)
	builder.Path(start, term, term)

	buildExpr(p, builder, "", start, term, term)
	buildExpr(p, builder, "term", term, term, term)

	builder.Path(quoted, term, term)
	builder.Path(quoted, parenClose, parenClose)

	builder.Path(number, term, term)

	builder.Path(ltrue, term, term)

	builder.Path(lfalse, term, term)

	builder.Path(term, number, number)
	builder.Path(term, quoted, quoted)
	builder.Path(term, ltrue, ltrue)
	builder.Path(term, lfalse, lfalse)

	builder.Path(let, identifier, "let-identifier")
	builder.Path("let-identifier", identifier, "let-type-identifier")
	builder.Path("let-type-identifier", equals, "let-equals")
	buildExpr(p, builder, "let", "let-equals", term, term)
	builder.WhenEntering("let-identifier", p.createIdentifier)
	builder.WhenEntering("let-type-identifier", p.createIdentifier)

	builder.WhenEntering(quoted, p.createStringLiteral)
	builder.WhenEntering(parenClose, p.closeNode)
	builder.WhenEntering(term, p.closeNode)
	builder.WhenEntering(number, p.createNumberLiteral)
	builder.WhenEntering(ltrue, p.createBooleanLiteral)
	builder.WhenEntering(lfalse, p.createBooleanLiteral)
	builder.WhenEntering(operator, p.createOperator)
	builder.WhenEntering(let, p.createLet)

	builder.Accept(term)

	return builder.Start(start)
}

// Builds rules for when expressions are allowed.
// This creates a new section of the DFA with a prefix
// that is entered following a particular token.
// For example, the expression allowed after the assignment
// operator. These nodes are in the let-specific section of the
// the DFA.
func buildExpr(p *parser, b dfa.MachineBuilder, prefix string, from string, returnVia string, returnTo string) {
	if len(prefix) > 0 {
		prefix = prefix + "-expr-"
	} else {
		prefix = "expr-"
	}

	exprNumber := prefix + lex.LNumber.String()
	exprString := prefix + lex.LQuoted.String()
	exprBoolTrue := prefix + lex.LBoolTrue.String()
	exprBoolFalse := prefix + lex.LBoolFalse.String()
	exprOperator := prefix + lex.LOperator.String()
	exprIdentifier := prefix + lex.LIdentifier.String()
	exprParenOpen := prefix + lex.LParenOpen.String()

	b.Path(from, number, exprNumber)
	b.Path(from, identifier, exprIdentifier)
	b.Path(from, ltrue, exprBoolTrue)
	b.Path(from, lfalse, exprBoolFalse)
	b.Path(exprIdentifier, operator, exprOperator)
	b.Path(exprIdentifier, parenOpen, exprParenOpen)
	b.Path(exprIdentifier, parenClose, from)
	b.Path(exprIdentifier, returnVia, returnTo)
	b.Path(exprParenOpen, quoted, exprString)
	b.Path(exprParenOpen, number, exprNumber)
	b.Path(exprParenOpen, identifier, exprIdentifier)
	b.Path(exprOperator, number, exprNumber)
	b.Path(exprOperator, identifier, exprIdentifier)
	b.Path(exprOperator, quoted, exprString)
	b.Path(exprNumber, operator, exprOperator)
	b.Path(exprNumber, returnVia, returnTo)
	b.Path(exprNumber, parenClose, from)
	b.Path(exprString, parenClose, from)
	b.Path(exprString, operator, exprOperator)
	b.Path(exprBoolTrue, returnVia, returnTo)
	b.Path(exprBoolFalse, returnVia, returnTo)

	b.WhenEntering(exprNumber, p.createNumberLiteral)
	b.WhenEntering(exprString, p.createStringLiteral)
	b.WhenEntering(exprBoolTrue, p.createBooleanLiteral)
	b.WhenEntering(exprBoolFalse, p.createBooleanLiteral)
	b.WhenEntering(exprIdentifier, p.createIdentifier)
	b.WhenEntering(exprOperator, p.createOperator)
}
