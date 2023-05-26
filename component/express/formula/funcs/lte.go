// Author: huaxr
// Time:   2021/7/7 下午2:42
// Git:    huaxr

package funcs

import (
	"github.com/huaxr/magicflow/component/logx"
	"github.com/spf13/cast"
	"github.com/yidane/formula/opt"
	"reflect"
)

type LteFunction struct {
}

func (*LteFunction) Name() string {
	return "lte"
}

func (f *LteFunction) Evaluate(context *opt.FormulaContext, args ...*opt.LogicalExpression) (*opt.Argument, error) {
	err := opt.MatchTwoArgument(f.Name(), args...)
	if err != nil {
		logx.L().Errorf("LteFunction err %v", err)
		return nil, err
	}

	val, err := (*args[0]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("LteFunction err %v", err)
		return nil, err
	}

	val2, err := (*args[1]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("LteFunction err %v", err)
		return nil, err
	}

	v1 := cast.ToFloat64(cast.ToString(val))
	v2 := cast.ToFloat64(cast.ToString(val2))
	return opt.NewArgumentWithType(v1 <= v2, reflect.Bool), nil
}
