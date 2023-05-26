// Author: huaxr
// Time:   2021/7/7 下午2:43
// Git:    huaxr

package funcs

import (
	"errors"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/yidane/formula/opt"
	"reflect"
)

type ContainsFunction struct {
}

func (*ContainsFunction) Name() string {
	return "contains"
}

func containsList(val []interface{}, v interface{}) bool {
	for _, i := range val {
		if i == v {
			return true
		}
	}
	return false
}

func containsMap(val map[string]interface{}, v interface{}) bool {
	for k, _ := range val {
		if k == v {
			return true
		}
	}
	return false
}

func (f *ContainsFunction) Evaluate(context *opt.FormulaContext, args ...*opt.LogicalExpression) (*opt.Argument, error) {
	err := opt.MatchTwoArgument(f.Name(), args...)
	if err != nil {
		logx.L().Errorf("ContainsFunction err %v", err)
		return nil, err
	}

	val, err := (*args[0]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("ContainsFunction %v", err)
		return nil, err
	}

	val2, err := (*args[1]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("ContainsFunction %v", err)
		return nil, err
	}

	switch reflect.TypeOf(val.Value).Kind() {
	case reflect.Array, reflect.Slice:
		return opt.NewArgumentWithType(containsList(val.Value.([]interface{}), val2), reflect.Bool), nil
	case reflect.Map:
		return opt.NewArgumentWithType(containsMap(val.Value.(map[string]interface{}), val2), reflect.Bool), nil
	default:
		return nil, errors.New("hasKey function first parameter must be slice or map")
	}
}
