// Generated from Express.g4 by ANTLR 4.7.

package antlr // Express

import "github.com/antlr/antlr4/runtime/Go/antlr"

// ExpressListener is a complete listener for a parse tree produced by ExpressParser.
type ExpressListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterEQUALS is called when entering the EQUALS production.
	EnterEQUALS(c *EQUALSContext)

	// EnterEXP is called when entering the EXP production.
	EnterEXP(c *EXPContext)

	// EnterNUMBER is called when entering the NUMBER production.
	EnterNUMBER(c *NUMBERContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitEQUALS is called when exiting the EQUALS production.
	ExitEQUALS(c *EQUALSContext)

	// ExitEXP is called when exiting the EXP production.
	ExitEXP(c *EXPContext)

	// ExitNUMBER is called when exiting the NUMBER production.
	ExitNUMBER(c *NUMBERContext)
}
