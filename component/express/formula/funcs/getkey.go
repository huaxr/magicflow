package funcs

import (
	"errors"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/spf13/cast"
	"github.com/yidane/formula/opt"
	"reflect"
)

type GetKeyFunction struct {
}

func (*GetKeyFunction) Name() string {
	return "getKey"
}

func (f *GetKeyFunction) Evaluate(context *opt.FormulaContext, args ...*opt.LogicalExpression) (*opt.Argument, error) {
	err := opt.MatchTwoArgument(f.Name(), args...)
	if err != nil {
		logx.L().Errorf("GetKeyFunction err %v", err)
		return nil, err
	}

	val, err := (*args[0]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("GetKeyFunction err %v", err)
		return nil, err
	}

	val2, err := (*args[1]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("GetKeyFunction err %v", err)
		return nil, err
	}

	index := cast.ToString(val2)

	switch reflect.TypeOf(val.Value).Kind() {
	case reflect.Map:
		return opt.NewArgumentWithType(val.Value.(map[string]interface{})[index], reflect.Interface), nil
	default:
		return nil, errors.New("getKey function first parameter must be map")
	}
}
