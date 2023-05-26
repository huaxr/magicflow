// Author: huaxr
// Time:   2021/6/25 上午10:26
// Git:    huaxr

package consensus

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/coreos/etcd/clientv3"
	"github.com/samuel/go-zookeeper/zk"
	"go.uber.org/atomic"
	"k8s.io/apimachinery/pkg/util/wait"
)

type KVType int

const (
	_ KVType = iota
	ETCD
	Zookeeper
)

var (
	once      sync.Once
	consensus Consensus
)

func LaunchIdGenerate(ctx context.Context, typ KVType) {
	NewConsensus(ctx, typ)
	InitIdGenerate(ctx)
}

func LaunchCampaign(ctx context.Context, typ KVType) {
	cli := NewConsensus(ctx, typ)
	go wait.UntilWithContext(ctx, cli.Elect, time.Second*5)
}

func NewConsensus(ctx context.Context, typ KVType) Consensus {
	if consensus != nil {
		return consensus
	}
	switch typ {
	case ETCD:
		addresses := confutil.GetConf().Dcc.Hosts
		logx.L().Infof("connect etcd addresses: %v", addresses)
		once.Do(func() {
			ctx, _ := context.WithCancel(ctx)
			cli, err := clientv3.New(clientv3.Config{
				Endpoints:   strings.Split(addresses, ","),
				DialTimeout: 5 * time.Second,
			})
			if err != nil {
				panic(err)
			}
			logx.L().Infof("etcd connect success")
			econsensus := new(EtcdConsensus)
			econsensus.client = cli
			econsensus.ctx = ctx
			econsensus.stop = atomic.NewBool(false)
			consensus = econsensus
		})

	case Zookeeper:
		panic("not imp")
		addresses := confutil.GetConf().Configuration.None
		logx.L().Infof("connect zk addresses: %v", addresses)
		once.Do(func() {
			ctx, _ := context.WithCancel(ctx)
			eventCallbackOption := zk.WithEventCallback(watchcallback)
			conn, _, err := zk.Connect(
				strings.Split(addresses, ","),
				time.Second*5,
				eventCallbackOption)
			if err != nil {
				panic(err)
			}
			logx.L().Infof("zookeeper connect success")
			zconsensus = new(ZkConsensus)
			zconsensus.client = conn
			zconsensus.ctx = ctx
			zconsensus.stop = atomic.NewBool(false)

			//consensus = zconsensus
		})

	default:
		panic("not implement yet")
	}

	return consensus
}
