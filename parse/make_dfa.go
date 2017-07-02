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
var comma = lex.LComma.String()

func buildDfa(p *parser) (dfa.Machine, error) {

	builder := dfa.NewMachineBuilder()

	builder.Path(start, quoted, quoted)
	builder.Path(start, ltrue, ltrue)
	builder.Path(start, lfalse, lfalse)
	builder.Path(start, let, let)
	builder.Path(start, term, start)

	defaultExprPrefix := buildExpr(p, builder, "", start, term, start)

	builder.Path(defaultExprPrefix+identifier, equals, equals)
	buildExpr(p, builder, "assignment", equals, term, start)

	builder.Path(quoted, term, start)
	builder.Path(quoted, parenClose, parenClose)

	builder.Path(number, term, start)

	builder.Path(ltrue, term, start)

	builder.Path(lfalse, term, start)

	builder.Path(let, identifier, "let-identifier")
	builder.Path("let-identifier", identifier, "let-type-identifier")
	builder.Path("let-type-identifier", term, start)
	builder.Path("let-type-identifier", equals, "let-equals")
	buildExpr(p, builder, "let", "let-equals", term, start)
	builder.WhenEntering("let-identifier", p.createIdentifier)
	builder.WhenEntering("let-type-identifier", p.createIdentifier)

	builder.WhenEntering(quoted, p.createStringLiteral)
	builder.WhenEntering(parenClose, p.closeGroupOrFunction)
	builder.WhenEntering(start, p.closeStatement)
	builder.WhenEntering(number, p.createNumberLiteral)
	builder.WhenEntering(ltrue, p.createBooleanLiteral)
	builder.WhenEntering(lfalse, p.createBooleanLiteral)
	builder.WhenEntering(operator, p.createOperator)
	builder.WhenEntering(let, p.createLet)
	builder.WhenEntering(parenOpen, p.createGroup)
	builder.WhenEntering(equals, p.createAssignment)

	builder.Accept(start)

	return builder.Start(start)
}

// Builds rules for when expressions are allowed.
// This creates a new section of the DFA with a prefix
// that is entered following a particular token.
// For example, the expression allowed after the assignment
// operator. These nodes are in the let-specific section of the
// the DFA.
//
// This function returns the prefix used so that the built portion
// of the AST can be extended externally.
// For example, at the beginning of statement we want to extend
// identifier to allow an equals to follow (for assignment).
func buildExpr(p *parser, b dfa.MachineBuilder, prefix string, from string, returnVia string, returnTo string) string {
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
	exprParenClose := prefix + lex.LParenClose.String()
	exprComma := prefix + lex.LComma.String()

	b.Path(from, number, exprNumber)
	b.Path(from, identifier, exprIdentifier)
	b.Path(from, quoted, exprString)
	b.Path(from, ltrue, exprBoolTrue)
	b.Path(from, lfalse, exprBoolFalse)
	b.Path(from, operator, exprOperator)
	b.Path(from, parenClose, exprParenClose)
	b.Path(exprIdentifier, operator, exprOperator)
	b.Path(exprIdentifier, parenOpen, exprParenOpen)
	b.Path(exprIdentifier, parenClose, exprParenClose)
	b.Path(exprIdentifier, returnVia, returnTo)
	b.Path(exprParenOpen, quoted, exprString)
	b.Path(exprParenOpen, number, exprNumber)
	b.Path(exprParenOpen, identifier, exprIdentifier)
	b.Path(exprParenOpen, parenOpen, exprParenOpen)
	b.Path(exprOperator, number, exprNumber)
	b.Path(exprOperator, identifier, exprIdentifier)
	b.Path(exprOperator, quoted, exprString)
	b.Path(exprOperator, ltrue, exprBoolTrue)
	b.Path(exprOperator, lfalse, exprBoolFalse)
	b.Path(exprOperator, parenOpen, exprParenOpen)
	b.Path(exprNumber, operator, exprOperator)
	b.Path(exprNumber, returnVia, returnTo)
	b.Path(exprNumber, parenClose, exprParenClose)
	b.Path(exprString, parenClose, exprParenClose)
	b.Path(exprString, operator, exprOperator)
	b.Path(exprString, comma, exprComma)
	b.Path(exprString, returnVia, returnTo)
	b.Path(exprComma, quoted, exprString)
	b.Path(exprBoolTrue, returnVia, returnTo)
	b.Path(exprBoolFalse, returnVia, returnTo)
	b.Path(exprBoolTrue, operator, exprOperator)
	b.Path(exprBoolFalse, returnVia, returnTo)
	b.Path(exprBoolTrue, parenClose, exprParenClose)
	b.Path(exprBoolFalse, parenClose, exprParenClose)
	b.Path(exprParenClose, operator, exprOperator)
	b.Path(exprParenClose, parenClose, exprParenClose)
	b.Path(exprParenClose, returnVia, returnTo)

	b.WhenEntering(exprNumber, p.createNumberLiteral)
	b.WhenEntering(exprString, p.createStringLiteral)
	b.WhenEntering(exprBoolTrue, p.createBooleanLiteral)
	b.WhenEntering(exprBoolFalse, p.createBooleanLiteral)
	b.WhenEntering(exprIdentifier, p.createIdentifier)
	b.WhenEntering(exprOperator, p.createOperator)
	b.WhenEntering(exprComma, p.closeArgument)
	b.WhenEntering(exprParenOpen, p.createGroup)
	b.WhenEntering(exprParenClose, p.closeGroupOrFunction)

	return prefix
}
