// Author: huaxr
// Time:   2021/7/6 下午5:20
// Git:    huaxr

package formula

import (
	"errors"
	"fmt"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/yidane/formula"
)

type Token string

const (
	And Token = " and "
	Or  Token = " or "
	// express could not contain this flag label.
	Default = "^x^"
)

type Expression struct {
	Condition string
	Params    []interface{}
}

// define the interfalce implement
type FormulaExpress struct {
	Exps  []*Expression
	Token Token
}

// Formula is an application scenarios of AntLr.
// we abbreviate the lexer g4 file to handle our expression instead of
// using this third library to achieve some custom functionality.
// @condition: equals([param0], "hello world"), paramIndex just traverse
// the @params slice to take over, which is $./$$./$$$. to fetch in the
// message context
func (d *FormulaExpress) EExecute() (bool, error) {
	var resBool = make([]bool, 0)
	var err error
	for _, exp := range d.Exps {
		expression := formula.NewExpression(exp.Condition)
		for index, param := range exp.Params {
			err = expression.AddParameter(fmt.Sprintf("param%v", index), param)
			if err != nil {
				logx.L().Errorf("ExecuteByFormula.err", "err when add parameter: %v", err)
				return false, err
			}
		}
		result, err := expression.Evaluate()
		if err != nil {
			logx.L().Errorf("ExecuteByFormula.err", "err when evaluate: %v", err)
			return false, err
		}

		v, ok := result.Value.(bool)
		if !ok {
			return false, errors.New("not boolean type")
		}
		resBool = append(resBool, v)
	}
	var v bool
	switch d.Token {
	case And:
		v = true
		for _, i := range resBool {
			v = v && i
		}
	case Or:
		for _, i := range resBool {
			v = v || i
		}
	case Default:
		v = resBool[0]
	default:
		panic("not implement")
	}
	return v, nil
}
