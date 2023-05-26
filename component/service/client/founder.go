// Author: XinRui Hua
// Time:   2022/3/21 下午6:22
// Git:    huaxr

package client

import (
	"context"

	"github.com/huaxr/magicflow/component/consensus"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/coreos/etcd/clientv3"
)

func ServerFound() {
	ctx := context.Background()
	cli := consensus.NewConsensus(ctx, consensus.ETCD)
	prefix := confutil.GetConf().Configuration.ServicesPrefix
	go watcher(cli.WatchKey(prefix))
	if res, err := cli.GetClient().Get(ctx, prefix, clientv3.WithPrefix()); err == nil {
		for _, i := range res.Kvs {
			rpcServices.register(string(i.Value))
		}
	} else {
		logx.L().Errorf("query services fail: %v", err)
	}
}
