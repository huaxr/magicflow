// Author: XinRui Hua
// Time:   2022/4/21 下午3:32
// Git:    huaxr

package promethu

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/ticker"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cast"

	"github.com/huaxr/magicflow/component/monitor/promethu/metric"
	"github.com/huaxr/magicflow/component/monitor/promethu/vars"
	"github.com/prometheus/client_golang/prometheus"
)

type manager struct {
	ctx context.Context
	// using tag
	tagDesc *prometheus.Desc
	// trace gauge
	traceDesc *prometheus.Desc
	// sync gauge
	syncDesc *prometheus.Desc

	ackDesc *prometheus.Desc

	// trigger plus
	// complete sub
	taskDesc      *prometheus.Desc
	summaryDesc   *prometheus.Desc
	histogramDesc *prometheus.Desc
}

func (c *manager) Name() string { return "prom_push" }
func (c *manager) Duration() *time.Ticker {
	return time.NewTicker(time.Duration(cast.ToInt(rand.Float64()*1000))*time.Millisecond + 3*time.Second)
}
func (c *manager) Heartbeat() {
	job := ticker.NewJob(c.ctx, c.Name(), c.Duration(), c.ticker)
	ticker.GetManager().Register(job)
}

func (c *manager) ticker() {
	if err := push.New(confutil.GetProm().GetPushGateWay(), "pushgateway").
		Collector(c).
		Grouping("host", fmt.Sprintf("%s:%s", toolutil.GetIp(), confutil.GetConf().Port.Service)). // add labels Groupings
		Push(); err != nil {
		logx.L().Errorf("Could not push completion time to Pushgateway:%v", err)
	}
}

// Describe simply sends the two Descs in the struct to the channel.
func (c *manager) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.tagDesc
	ch <- c.traceDesc
	ch <- c.syncDesc
	ch <- c.taskDesc
	ch <- c.ackDesc
}

func (c *manager) Collect(ch chan<- prometheus.Metric) {
	for _, i := range metric.GetMetric() {
		ch <- prometheus.MustNewConstMetric(
			c.tagDesc,
			prometheus.CounterValue,
			float64(i.GetCount()),
			i.GetKey(),
		)
		metric.Clean(i.GetTag())
	}

	// counter gauge
	ch <- prometheus.MustNewConstMetric(c.traceDesc, prometheus.GaugeValue, float64(vars.GetTraceCount()))
	ch <- prometheus.MustNewConstMetric(c.syncDesc, prometheus.GaugeValue, float64(vars.GetSyncCount()))
	ch <- prometheus.MustNewConstMetric(c.ackDesc, prometheus.GaugeValue, float64(vars.GetAckCount()))

	// pie chart
	ch <- prometheus.MustNewConstMetric(c.taskDesc, prometheus.GaugeValue, float64(vars.GetResidualTask()))
}

func (c *manager) collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstSummary(
		c.summaryDesc,
		4711, 403.34,
		map[float64]float64{0.5: 42.3, 0.9: 323.3},
		"200", "get",
	)

	ch <- prometheus.MustNewConstHistogram(
		c.histogramDesc,
		4711, 403.34,
		map[float64]uint64{25: 121, 50: 2403, 100: 3221, 200: 4233},
		"200", "get",
	)
}

func NewManager(ctx context.Context) *manager {
	return &manager{
		ctx: ctx,
		tagDesc: prometheus.NewDesc(
			fmt.Sprintf("flow_tag"),
			"tag distribute",
			[]string{"tag"}, // add keys here
			prometheus.Labels{},
		),

		traceDesc: prometheus.NewDesc(
			"trace_count",
			"trace key count",
			nil, nil),

		syncDesc: prometheus.NewDesc(
			"sync_count",
			"sync key count",
			nil, nil),

		ackDesc: prometheus.NewDesc(
			"ack_count",
			"acl key count",
			nil, nil),

		taskDesc: prometheus.NewDesc(
			"residual_task_count",
			"residual task count",
			// with app name
			[]string{}, nil),

		summaryDesc: prometheus.NewDesc(
			"summary_duration_seconds",
			"summary",
			[]string{"code", "method"},
			prometheus.Labels{},
		),
		histogramDesc: prometheus.NewDesc(
			"histogram_duration_seconds",
			"histogram",
			[]string{"code", "method"},
			prometheus.Labels{},
		),
	}
}
