// Author: huaxr
// Time:   2021/10/18 下午6:45
// Git:    huaxr

package antlr

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"log"
)

type expListener struct {
	*BaseExpressListener

	funcs  []string
	params []interface{}
	values []interface{}
}

func (e *expListener) ExitEQUALS(c *EQUALSContext) {
	log.Println("exit EQUASL", c.EQUALS().GetText())
}

func (e *expListener) ExitEXP(c *EXPContext) {
	log.Println("exit exp", c.EXP().GetText())
}

func (e *expListener) ExitNUMBER(c *NUMBERContext) {
	log.Println("exit number", c.NUMBER().GetText())
}

func handler(input string) {
	// Setup the input
	is := antlr.NewInputStream(input)

	// Create the Lexer
	lexer := NewExpressLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := NewExpressParser(stream)

	// Finally parse the expression (by walking the tree)
	var listener expListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, p.Start())

	log.Println(listener.funcs, listener.params, listener.values)
}
