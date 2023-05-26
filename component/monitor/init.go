// Author: XinRui Hua
// Time:   2022/4/15 下午6:17
// Git:    huaxr

package monitor

import (
	"context"
	"github.com/huaxr/magicflow/component/monitor/promethu"
)

// launch promethues by pull or push
// if type is Pull, it will listening on 8800 to handler /metricx request
// on the other hand, Push will push every metricx to the pushgateway in
// order to avoiding network connection error.
func LaunchMonitor(ctx context.Context, typ promethu.Reporter) {
	go promethu.LaunchProm(ctx, typ)
}
