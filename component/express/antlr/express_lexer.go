// Generated from Express.g4 by ANTLR 4.7.

package antlr

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 10, 101,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 3, 2, 6, 2, 21, 10, 2, 13, 2, 14, 2, 22,
	3, 3, 6, 3, 26, 10, 3, 13, 3, 14, 3, 27, 3, 4, 7, 4, 31, 10, 4, 12, 4,
	14, 4, 34, 11, 4, 3, 4, 3, 4, 7, 4, 38, 10, 4, 12, 4, 14, 4, 41, 11, 4,
	3, 5, 7, 5, 44, 10, 5, 12, 5, 14, 5, 47, 11, 5, 3, 5, 3, 5, 7, 5, 51, 10,
	5, 12, 5, 14, 5, 54, 11, 5, 3, 6, 7, 6, 57, 10, 6, 12, 6, 14, 6, 60, 11,
	6, 3, 6, 3, 6, 7, 6, 64, 10, 6, 12, 6, 14, 6, 67, 11, 6, 3, 7, 7, 7, 70,
	10, 7, 12, 7, 14, 7, 73, 11, 7, 3, 7, 3, 7, 7, 7, 77, 10, 7, 12, 7, 14,
	7, 80, 11, 7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3,
	8, 3, 8, 3, 8, 3, 8, 3, 9, 6, 9, 96, 10, 9, 13, 9, 14, 9, 97, 3, 9, 3,
	9, 2, 2, 10, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 10, 3, 2,
	5, 3, 2, 50, 59, 4, 2, 67, 92, 99, 124, 5, 2, 11, 12, 15, 15, 34, 34, 2,
	111, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2,
	2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2, 17, 3,
	2, 2, 2, 3, 20, 3, 2, 2, 2, 5, 25, 3, 2, 2, 2, 7, 32, 3, 2, 2, 2, 9, 45,
	3, 2, 2, 2, 11, 58, 3, 2, 2, 2, 13, 71, 3, 2, 2, 2, 15, 81, 3, 2, 2, 2,
	17, 95, 3, 2, 2, 2, 19, 21, 9, 2, 2, 2, 20, 19, 3, 2, 2, 2, 21, 22, 3,
	2, 2, 2, 22, 20, 3, 2, 2, 2, 22, 23, 3, 2, 2, 2, 23, 4, 3, 2, 2, 2, 24,
	26, 9, 3, 2, 2, 25, 24, 3, 2, 2, 2, 26, 27, 3, 2, 2, 2, 27, 25, 3, 2, 2,
	2, 27, 28, 3, 2, 2, 2, 28, 6, 3, 2, 2, 2, 29, 31, 7, 38, 2, 2, 30, 29,
	3, 2, 2, 2, 31, 34, 3, 2, 2, 2, 32, 30, 3, 2, 2, 2, 32, 33, 3, 2, 2, 2,
	33, 35, 3, 2, 2, 2, 34, 32, 3, 2, 2, 2, 35, 39, 7, 48, 2, 2, 36, 38, 5,
	5, 3, 2, 37, 36, 3, 2, 2, 2, 38, 41, 3, 2, 2, 2, 39, 37, 3, 2, 2, 2, 39,
	40, 3, 2, 2, 2, 40, 8, 3, 2, 2, 2, 41, 39, 3, 2, 2, 2, 42, 44, 7, 34, 2,
	2, 43, 42, 3, 2, 2, 2, 44, 47, 3, 2, 2, 2, 45, 43, 3, 2, 2, 2, 45, 46,
	3, 2, 2, 2, 46, 48, 3, 2, 2, 2, 47, 45, 3, 2, 2, 2, 48, 52, 7, 46, 2, 2,
	49, 51, 7, 34, 2, 2, 50, 49, 3, 2, 2, 2, 51, 54, 3, 2, 2, 2, 52, 50, 3,
	2, 2, 2, 52, 53, 3, 2, 2, 2, 53, 10, 3, 2, 2, 2, 54, 52, 3, 2, 2, 2, 55,
	57, 7, 34, 2, 2, 56, 55, 3, 2, 2, 2, 57, 60, 3, 2, 2, 2, 58, 56, 3, 2,
	2, 2, 58, 59, 3, 2, 2, 2, 59, 61, 3, 2, 2, 2, 60, 58, 3, 2, 2, 2, 61, 65,
	7, 42, 2, 2, 62, 64, 7, 34, 2, 2, 63, 62, 3, 2, 2, 2, 64, 67, 3, 2, 2,
	2, 65, 63, 3, 2, 2, 2, 65, 66, 3, 2, 2, 2, 66, 12, 3, 2, 2, 2, 67, 65,
	3, 2, 2, 2, 68, 70, 7, 34, 2, 2, 69, 68, 3, 2, 2, 2, 70, 73, 3, 2, 2, 2,
	71, 69, 3, 2, 2, 2, 71, 72, 3, 2, 2, 2, 72, 74, 3, 2, 2, 2, 73, 71, 3,
	2, 2, 2, 74, 78, 7, 43, 2, 2, 75, 77, 7, 34, 2, 2, 76, 75, 3, 2, 2, 2,
	77, 80, 3, 2, 2, 2, 78, 76, 3, 2, 2, 2, 78, 79, 3, 2, 2, 2, 79, 14, 3,
	2, 2, 2, 80, 78, 3, 2, 2, 2, 81, 82, 7, 103, 2, 2, 82, 83, 7, 115, 2, 2,
	83, 84, 7, 119, 2, 2, 84, 85, 7, 99, 2, 2, 85, 86, 7, 110, 2, 2, 86, 87,
	7, 117, 2, 2, 87, 88, 3, 2, 2, 2, 88, 89, 5, 11, 6, 2, 89, 90, 5, 7, 4,
	2, 90, 91, 5, 9, 5, 2, 91, 92, 5, 3, 2, 2, 92, 93, 5, 13, 7, 2, 93, 16,
	3, 2, 2, 2, 94, 96, 9, 4, 2, 2, 95, 94, 3, 2, 2, 2, 96, 97, 3, 2, 2, 2,
	97, 95, 3, 2, 2, 2, 97, 98, 3, 2, 2, 2, 98, 99, 3, 2, 2, 2, 99, 100, 8,
	9, 2, 2, 100, 18, 3, 2, 2, 2, 14, 2, 22, 27, 32, 39, 45, 52, 58, 65, 71,
	78, 97, 3, 8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames []string

var lexerSymbolicNames = []string{
	"", "NUMBER", "CHAR", "EXP", "Comma", "LBracket", "RBracket", "EQUALS",
	"WHITESPACE",
}

var lexerRuleNames = []string{
	"NUMBER", "CHAR", "EXP", "Comma", "LBracket", "RBracket", "EQUALS", "WHITESPACE",
}

type ExpressLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewExpressLexer(input antlr.CharStream) *ExpressLexer {

	l := new(ExpressLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "Express.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// ExpressLexer tokens.
const (
	ExpressLexerNUMBER     = 1
	ExpressLexerCHAR       = 2
	ExpressLexerEXP        = 3
	ExpressLexerComma      = 4
	ExpressLexerLBracket   = 5
	ExpressLexerRBracket   = 6
	ExpressLexerEQUALS     = 7
	ExpressLexerWHITESPACE = 8
)
