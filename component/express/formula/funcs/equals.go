// Author: huaxr
// Time:   2021/7/7 上午11:46
// Git:    huaxr

package funcs

import (
	"github.com/huaxr/magicflow/component/logx"
	"github.com/spf13/cast"
	"github.com/yidane/formula/opt"
	"reflect"
	"strings"
)

type EqualsFunction struct {
}

func (*EqualsFunction) Name() string {
	return "equals"
}

func (f *EqualsFunction) Evaluate(context *opt.FormulaContext, args ...*opt.LogicalExpression) (*opt.Argument, error) {
	err := opt.MatchTwoArgument(f.Name(), args...)
	if err != nil {
		logx.L().Errorf("EqualsFunction err %v", err)
		return nil, err
	}

	val, err := (*args[0]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("EqualsFunction err %v", err)
		return nil, err
	}

	val2, err := (*args[1]).Evaluate(context)
	if err != nil {
		logx.L().Errorf("EqualsFunction err %v", err)
		return nil, err
	}

	return opt.NewArgumentWithType(cast.ToString(val) == strings.Trim(cast.ToString(val2), "'"), reflect.Bool), nil
}
