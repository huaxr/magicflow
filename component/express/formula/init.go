// Author: huaxr
// Time:   2021/7/6 下午5:16
// Git:    huaxr

package formula

import (
	"github.com/huaxr/magicflow/component/express/formula/funcs"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/yidane/formula"
	"github.com/yidane/formula/opt"
	"sync"
)

var once sync.Once

func LaunchFormulaFuncs() {
	once.Do(func() {
		var functions = make([]opt.Function, 0)
		functions = append(functions,
			&funcs.LenFunction{},
			&funcs.EqualsFunction{},
			&funcs.GteFunction{},
			&funcs.GtFunction{},
			&funcs.LteFunction{},
			&funcs.LtFunction{},
			&funcs.ContainsFunction{},
			//&GetindexFunction{},
			&funcs.GetKeyFunction{},
		)

		for _, i := range functions {
			// wtf!!!
			M := i
			err := formula.Register(&M)
			if err != nil {
				logx.L().Errorf("err when launch func:%+v", i.Name())
			}
		}

	})
}
