// Author: huaxr
// Time: 2022/7/6 10:08 上午
// Git: huaxr

package ticker

import (
	"time"

	"github.com/huaxr/magicflow/component/logx"
)

type Tick interface {
	Name() string
	Duration() *time.Ticker
	Heartbeat()
}

func RegisterTick(t Tick) {
	logx.L().Infof("register ticker for %v", t.Name())
	t.Heartbeat()
}
