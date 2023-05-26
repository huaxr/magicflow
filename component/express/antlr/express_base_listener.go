// Generated from Express.g4 by ANTLR 4.7.

package antlr // Express

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseExpressListener is a complete listener for a parse tree produced by ExpressParser.
type BaseExpressListener struct{}

var _ ExpressListener = &BaseExpressListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExpressListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExpressListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExpressListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExpressListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseExpressListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseExpressListener) ExitStart(ctx *StartContext) {}

// EnterEQUALS is called when production EQUALS is entered.
func (s *BaseExpressListener) EnterEQUALS(ctx *EQUALSContext) {}

// ExitEQUALS is called when production EQUALS is exited.
func (s *BaseExpressListener) ExitEQUALS(ctx *EQUALSContext) {}

// EnterEXP is called when production EXP is entered.
func (s *BaseExpressListener) EnterEXP(ctx *EXPContext) {}

// ExitEXP is called when production EXP is exited.
func (s *BaseExpressListener) ExitEXP(ctx *EXPContext) {}

// EnterNUMBER is called when production NUMBER is entered.
func (s *BaseExpressListener) EnterNUMBER(ctx *NUMBERContext) {}

// ExitNUMBER is called when production NUMBER is exited.
func (s *BaseExpressListener) ExitNUMBER(ctx *NUMBERContext) {}
