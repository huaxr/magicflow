// Author: huaxr
// Time:   2021/7/6 下午5:08
// Git:    huaxr

package funcs

import (
	"errors"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/yidane/formula/opt"
	"reflect"
)

type LenFunction struct {
}

func (*LenFunction) Name() string {
	return "len"
}

func (f *LenFunction) Evaluate(context *opt.FormulaContext, args ...*opt.LogicalExpression) (*opt.Argument, error) {
	err := opt.MatchOneArgument(f.Name(), args...)
	if err != nil {
		logx.L().Errorf("EqualsFunction.err", "%v", err)
		return nil, err
	}

	val, err := (*args[0]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("LenFunction.err", "%v", err)
		return nil, err
	}

	switch reflect.TypeOf(val.Value).Kind() {
	case reflect.Array, reflect.Slice:
		return opt.NewArgumentWithType(len(val.Value.([]interface{})), reflect.Int), nil
	case reflect.Map:
		return opt.NewArgumentWithType(len(val.Value.(map[string]interface{})), reflect.Int), nil
	case reflect.String:
		return opt.NewArgumentWithType(len(val.Value.(string)), reflect.Int), nil
	default:
		return nil, errors.New("not support parameter, only allowed map,slice,string")
	}
}
