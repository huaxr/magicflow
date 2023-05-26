// Author: XinRui Hua
// Time:   2022/4/21 下午2:49
// Git:    huaxr

package promethu

import (
	"context"

	"github.com/huaxr/magicflow/component/ticker"

	"github.com/huaxr/magicflow/component/monitor/promethu/metric"
)

func PromPush(ctx context.Context) {
	metric.RegisterMetrics()
	worker := NewManager(ctx)
	ticker.RegisterTick(worker)
}
