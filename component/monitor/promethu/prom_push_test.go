// Author: XinRui Hua
// Time:   2022/4/25 上午10:58
// Git:    huaxr

package promethu

import (
	ctx2 "context"
	"math/rand"
	"testing"
	"time"

	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/ticker"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type testx struct {
	// using tag
	tDesc *prometheus.Desc
}

// Describe simply sends the two Descs in the struct to the channel.
func (c *testx) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.tDesc
}

func (c *testx) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.tDesc, prometheus.GaugeValue, rand.Float64()*100, toolutil.GetShortRandomString(1))
}

func TestPush(t *testing.T) {
	worker := &testx{tDesc: prometheus.NewDesc(
		"zzz",
		"xxx",
		[]string{"tmp"}, nil)}

	tick := time.NewTicker(1 * time.Second)
	job := ticker.NewJob(ctx2.Background(), "local_state_hb", tick, func() {
		if err := push.New("http://127.0.0.1:9091", "pushgateway").
			Collector(worker).
			Grouping("host", toolutil.GetShortRandomString(1)). // add labels Groupings
			Push(); err != nil {
			logx.L().Errorf("Could not push completion time to Pushgateway:%v", err)
		}
	})
	ticker.GetManager().Register(job)

}
