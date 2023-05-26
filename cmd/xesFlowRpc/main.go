package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huaxr/magicflow/component/consensus"
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/dcc"
	"github.com/huaxr/magicflow/component/dispatch"
	"github.com/huaxr/magicflow/component/helper/console"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/monitor"
	"github.com/huaxr/magicflow/component/monitor/healthy"
	"github.com/huaxr/magicflow/component/monitor/promethu"
	"github.com/huaxr/magicflow/component/service"
	"github.com/huaxr/magicflow/component/ticker"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/urfave/cli/v2"
)

var (
	exit        = os.Exit
	ctx, cancel = context.WithCancel(context.Background())
)

func flow(c *cli.Context) error {
	c.Context = ctx
	ticker.LaunchJobManager()
	orm.LaunchDbEngine()
	//kv.LunchRedis(ctx)
	dcc.LaunchDcc(ctx, consensus.ETCD)
	consensus.LaunchIdGenerate(ctx, consensus.ETCD)
	core.LaunchCore(ctx, dispatch.NSQ)

	if confutil.GetConf().Switch.EnableMonitor {
		monitor.LaunchMonitor(ctx, promethu.Pull)
	}

	if confutil.GetConf().Switch.EnableHealthyCheck {
		healthy.LaunchHealthy(ctx)
	}

	if confutil.GetConf().Switch.EnableMasterElect {
		consensus.LaunchCampaign(ctx, consensus.ETCD)
	}

	return service.LaunchRpcServer(ctx, confutil.GetConf().Port.Service)
}

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1)

	go func() {
		v := <-c
		logx.L().Warnf("exit signal %v received, system process abort.", v)
		cancel()
		time.Sleep(3 * time.Second)
		exit(1)
	}()

	if app := console.LaunchApp(flow); app.Run(os.Args) != nil {
		exit(1)
	}
}
