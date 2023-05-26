// Author: huaxr
// Time:   2021/9/18 上午11:59
// Git:    huaxr

package healthy

import (
	"context"
	"strings"
	"time"

	"github.com/huaxr/magicflow/component/ticker"

	"github.com/huaxr/magicflow/component/consensus"
	"github.com/huaxr/magicflow/component/dcc"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/ssrf"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/spf13/cast"
)

type channelStats struct{}

// return app to debug print the apps once in the log console.
func do(app string) string {
	var err error
	cli := ssrf.NewHttpClient(confutil.GetConf().Queue.Nsq.Admin)
	res := cli.ReportChannelStats(core.GetWorkerTopic(app))
	// when close app
	if res == nil {
		return ""
	}

	ns, ok := core.G().GetNamespace(app)
	// bypass TopicName blank cause consider this situation:
	// we delete this topic on nsqadmin platform which option
	// will close the app forever.
	if res.TopicName != "" {
		if !ok {
			return ""
		}

		if res.ClientCount == 0 && ns.Working() {
			logx.L().Infof("master send etcd to close app:%v", app)
			err = dcc.GetDcc().DccSwitchApp(app, dcc.CLOSE)
			if err != nil {
				logx.L().Errorf("channlestate do %v", err)
			}
		}
		if res.ClientCount > 0 && !ns.Working() {
			logx.L().Infof("master send etcd to open app:%v", app)
			err = dcc.GetDcc().DccSwitchApp(app, dcc.OPEN)
			if err != nil {
				logx.L().Errorf("channlestate do %v", err)
			}
		}
		return ""
	} else {
		if ns.Working() {
			logx.L().Infof("master send etcd to close app:%v", app)
			err = dcc.GetDcc().DccSwitchApp(app, dcc.CLOSE)
			if err != nil {
				logx.L().Errorf("channlestate do %v", err)
			}
		}
		return app
	}

}

func (c *channelStats) Name() string { return "channel_report" }
func (c *channelStats) Duration() *time.Ticker {
	t := cast.ToInt(strings.TrimRight(confutil.GetConf().Configuration.ChannelReportInterval, "s"))
	return time.NewTicker(time.Duration(t) * time.Second)
}
func (c *channelStats) Heartbeat() {
	c.ticker()
	job := ticker.NewJob(context.Background(), c.Name(), c.Duration(), c.ticker)
	ticker.GetManager().Register(job)
}

func (c *channelStats) ticker() {
	var diff []string
	for _, app := range core.G().GetAllNs() {
		// notify pods
		should := do(app)
		if should != "" {
			diff = append(diff, should)
		}
	}
	if len(diff) > 0 {
		//logx.L().Debugf("report diff: %v", diff)
	}
}

func (c *channelStats) Report() {
	if confutil.GetConf().IsLocalEnv() {
		logx.L().Infof("local env start channel stats monitoring")
		ticker.RegisterTick(c)
	} else {
		consensus.MasterFunc(c.Heartbeat)
	}
}

func (c *channelStats) String() string {
	return "channelStats"
}
