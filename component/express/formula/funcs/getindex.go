package funcs

import (
	"errors"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/spf13/cast"
	"github.com/yidane/formula/opt"
	"reflect"
)

type GetIndexFunction struct {
}

func (*GetIndexFunction) Name() string {
	return "getIndex"
}

func (f *GetIndexFunction) Evaluate(context *opt.FormulaContext, args ...*opt.LogicalExpression) (*opt.Argument, error) {
	err := opt.MatchTwoArgument(f.Name(), args...)
	if err != nil {
		logx.L().Errorf("EqualsFunction.err", "%v", err)
		return nil, err
	}

	val, err := (*args[0]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("GetIndexFunction.err", "%v", err)
		return nil, err
	}

	val2, err := (*args[1]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("GetIndexFunction.err", "%v", err)
		return nil, err
	}

	index := cast.ToInt(val2)

	switch reflect.TypeOf(val.Value).Kind() {
	case reflect.Array, reflect.Slice:
		return opt.NewArgumentWithType(val.Value.([]interface{})[index], reflect.Int), nil
	default:
		return nil, errors.New("getIndex function first parameter must be array or slice")
	}
}
