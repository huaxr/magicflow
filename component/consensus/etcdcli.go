package consensus

import (
	"context"

	"github.com/huaxr/magicflow/pkg/confutil"

	"time"

	"github.com/huaxr/magicflow/component/logx"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"go.uber.org/atomic"
)

type EtcdConsensus struct {
	ctx    context.Context
	client *clientv3.Client
	stop   *atomic.Bool
}

func (c *EtcdConsensus) Close() {
	c.stop.Store(true)
	c.client.Close()
}

func (c *EtcdConsensus) Put(key, val string) (err error) {
	_, err = c.client.Put(c.ctx, key, val)
	if err != nil {
		logx.L().Errorf("consensus put err :%v", err)
	}
	return
}

func (c *EtcdConsensus) Get(key string) [][]byte {
	var res = make([][]byte, 0)
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	resp, err := c.client.Get(ctx, key, clientv3.WithPrefix())
	defer cancel()
	if err != nil {
		logx.L().Errorf("consensus get err :%v", err)
		return nil
	}
	for _, ev := range resp.Kvs {
		logx.L().Infof("consensus success :%v, %v", ev.Key, ev.Value)
		res = append(res, ev.Value)
	}

	return res
}

func (c *EtcdConsensus) GetClient() *clientv3.Client {
	return c.client
}

func (c *EtcdConsensus) WatchKey(key string) clientv3.WatchChan {
	return c.client.Watch(c.ctx, key, clientv3.WithPrefix())
}

func (c *EtcdConsensus) Delete(key string) (resp *clientv3.DeleteResponse, err error) {
	return c.client.Delete(c.ctx, key, clientv3.WithPrefix())
}

func (c *EtcdConsensus) Elect(ctx context.Context) {
	session, err := concurrency.NewSession(consensus.GetClient(), concurrency.WithTTL(8))
	if err != nil {
		// k8s alarm here. why?
		logx.L().Warnf("elect concurrency.NewSession err: %v", err.Error())
		return
	}
	e := concurrency.NewElection(session, confutil.GetConf().Configuration.ElectionPrefix)

	// Campaign puts a value as eligible for the election. It blocks until
	// it is elected, an error occurs, or the context is cancelled.
	if err = e.Campaign(context.TODO(), leader); err != nil {
		// k8s alarm here. why?
		logx.L().Warnf("elect campaign err:%v", err)
		return
	}

	logx.L().Infof("Elect success")

	leaderFlag.Store(true)

	select {
	case <-session.Done():
		leaderFlag.Store(false)
		logx.L().Warnf("elect expired restart elect")
	}
}
