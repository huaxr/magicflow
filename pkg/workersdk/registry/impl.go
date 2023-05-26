// Author: huaxr
// Time:   2021/9/3 下午12:54
// Git:    huaxr

package registry

import (
	"github.com/huaxr/magicflow/pkg/workersdk/flow"
)

type Task interface {
	//fmt.Stringer
	Do(ctx flow.Context, input interface{}) (output interface{}, err error)
}

// register command
type Command interface {
	Execute(ctx flow.Context, params map[string]interface{})
}
