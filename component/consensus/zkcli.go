// Author: XinRui Hua
// Time:   2022/2/9 下午7:31
// Git:    huaxr

package consensus

import (
	"context"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/samuel/go-zookeeper/zk"
	"go.uber.org/atomic"
)

var (
	zconsensus *ZkConsensus
)

type ZkConsensus struct {
	ctx    context.Context
	client *zk.Conn
	stop   *atomic.Bool
}

var (
	flags int32 = zk.FlagEphemeral
	acls        = zk.WorldACL(zk.PermAll)
)

func watchcallback(event zk.Event) {
	logx.L().Infof("%+v", event)
}

func (z *ZkConsensus) Close() {
	z.stop.Store(true)
	z.client.Close()
}

func (z *ZkConsensus) Put(key, val string) (err error) {
	_, err = z.client.Create(key, toolutil.String2Byte(val), flags, acls)
	if err != nil {
		logx.L().Errorf("create err:%v", err)
		return err
	}

	return nil
}

func (z *ZkConsensus) Get(key string) [][]byte {
	return nil
}

func (z *ZkConsensus) del(path string) {
	_, stat, _ := z.client.Get(path)
	err := z.client.Delete(path, stat.Version)
	if err != nil {
		logx.L().Errorf("del err:%v", err)
		return
	}
}

func (z *ZkConsensus) watch(path string) {
	_, stat, _ := z.client.Get(path)
	err := z.client.Delete(path, stat.Version)
	if err != nil {
		logx.L().Errorf("del err:%v", err)
		return
	}
}
