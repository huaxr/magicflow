// Author: huaxr
// Time:   2021/12/27 上午11:29
// Git:    huaxr

package valuate

import (
	"errors"

	"github.com/Knetic/govaluate"
)

type ValuateExpress struct {
	Condition string
	Params    map[string]interface{}
}

func (v *ValuateExpress) EExecute() (bool, error) {
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(v.Condition, functions)
	if err != nil {
		return false, err
	}
	result, err := expression.Evaluate(v.Params)
	if err != nil {
		return false, err
	}
	_, ok := result.(bool)
	if !ok {
		return false, errors.New("not bool express")
	}
	return result.(bool), nil
}
