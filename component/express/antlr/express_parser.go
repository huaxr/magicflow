// Generated from Express.g4 by ANTLR 4.7.

package antlr // Express

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 10, 15, 4,
	2, 9, 2, 4, 3, 9, 3, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 5, 3, 13, 10,
	3, 3, 3, 2, 2, 4, 2, 4, 2, 2, 2, 14, 2, 6, 3, 2, 2, 2, 4, 12, 3, 2, 2,
	2, 6, 7, 5, 4, 3, 2, 7, 8, 7, 2, 2, 3, 8, 3, 3, 2, 2, 2, 9, 13, 7, 9, 2,
	2, 10, 13, 7, 5, 2, 2, 11, 13, 7, 3, 2, 2, 12, 9, 3, 2, 2, 2, 12, 10, 3,
	2, 2, 2, 12, 11, 3, 2, 2, 2, 13, 5, 3, 2, 2, 2, 3, 12,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames []string

var symbolicNames = []string{
	"", "NUMBER", "CHAR", "EXP", "Comma", "LBracket", "RBracket", "EQUALS",
	"WHITESPACE",
}

var ruleNames = []string{
	"start", "expression",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type ExpressParser struct {
	*antlr.BaseParser
}

func NewExpressParser(input antlr.TokenStream) *ExpressParser {
	this := new(ExpressParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Express.g4"

	return this
}

// ExpressParser tokens.
const (
	ExpressParserEOF        = antlr.TokenEOF
	ExpressParserNUMBER     = 1
	ExpressParserCHAR       = 2
	ExpressParserEXP        = 3
	ExpressParserComma      = 4
	ExpressParserLBracket   = 5
	ExpressParserRBracket   = 6
	ExpressParserEQUALS     = 7
	ExpressParserWHITESPACE = 8
)

// ExpressParser rules.
const (
	ExpressParserRULE_start      = 0
	ExpressParserRULE_expression = 1
)

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExpressParserRULE_start
	return p
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExpressParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *StartContext) EOF() antlr.TerminalNode {
	return s.GetToken(ExpressParserEOF, 0)
}

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExpressListener); ok {
		listenerT.EnterStart(s)
	}
}

func (s *StartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExpressListener); ok {
		listenerT.ExitStart(s)
	}
}

func (p *ExpressParser) Start() (localctx IStartContext) {
	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, ExpressParserRULE_start)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(4)
		p.Expression()
	}
	{
		p.SetState(5)
		p.Match(ExpressParserEOF)
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ExpressParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExpressParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) CopyFrom(ctx *ExpressionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type NUMBERContext struct {
	*ExpressionContext
}

func NewNUMBERContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NUMBERContext {
	var p = new(NUMBERContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *NUMBERContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NUMBERContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(ExpressParserNUMBER, 0)
}

func (s *NUMBERContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExpressListener); ok {
		listenerT.EnterNUMBER(s)
	}
}

func (s *NUMBERContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExpressListener); ok {
		listenerT.ExitNUMBER(s)
	}
}

type EQUALSContext struct {
	*ExpressionContext
}

func NewEQUALSContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EQUALSContext {
	var p = new(EQUALSContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *EQUALSContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EQUALSContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(ExpressParserEQUALS, 0)
}

func (s *EQUALSContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExpressListener); ok {
		listenerT.EnterEQUALS(s)
	}
}

func (s *EQUALSContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExpressListener); ok {
		listenerT.ExitEQUALS(s)
	}
}

type EXPContext struct {
	*ExpressionContext
}

func NewEXPContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EXPContext {
	var p = new(EXPContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *EXPContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EXPContext) EXP() antlr.TerminalNode {
	return s.GetToken(ExpressParserEXP, 0)
}

func (s *EXPContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExpressListener); ok {
		listenerT.EnterEXP(s)
	}
}

func (s *EXPContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExpressListener); ok {
		listenerT.ExitEXP(s)
	}
}

func (p *ExpressParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, ExpressParserRULE_expression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(10)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ExpressParserEQUALS:
		localctx = NewEQUALSContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(7)
			p.Match(ExpressParserEQUALS)
		}

	case ExpressParserEXP:
		localctx = NewEXPContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(8)
			p.Match(ExpressParserEXP)
		}

	case ExpressParserNUMBER:
		localctx = NewNUMBERContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(9)
			p.Match(ExpressParserNUMBER)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}
