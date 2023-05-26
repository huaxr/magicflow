// Author: XinRui Hua
// Time:   2022/3/21 下午5:41
// Git:    huaxr

package consensus

import (
	"context"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/coreos/etcd/clientv3"
	"strings"
	"time"
)

type KeepAlive struct {
	Key string
	Ttl int64
}

// rent key alive priod.
func (k KeepAlive) Alive(ctx context.Context) {
	sps := strings.Split(k.Key, "/")
	if len(sps) < 2 {
		logx.L().Errorf("key format err")
		return
	}

	client := NewConsensus(ctx, ETCD).GetClient()
	lease := clientv3.NewLease(client)
	// Apply for a ttl lease
	leaseResp, err := lease.Grant(ctx, k.Ttl)
	if err != nil {
		logx.L().Errorf("set lease fail:%v", err)
		return
	}
	// leaseid
	leaseID := leaseResp.ID
	// lease automatically
	leaseRespChan, err := lease.KeepAlive(ctx, leaseID)

	if err != nil {
		logx.L().Errorf("lease keepalive fail:%v", err)
		return
	}

	kv := clientv3.NewKV(client)

	_, err = kv.Put(context.TODO(), k.Key, sps[len(sps)-1], clientv3.WithLease(leaseID))
	if err != nil {
		logx.L().Errorf("put key fail:%v", err)
		return
	}
	logx.L().Infof("register put key %v success", k.Key)

	for {
		select {
		case leaseKeepResp := <-leaseRespChan:
			if leaseKeepResp == nil {
				// lease failure
				logx.L().Warnf("lease resp nil")
				return
			} else {
				// reply success
				time.Sleep(time.Duration(k.Ttl/2) * time.Second)
			}
		case <-ctx.Done():
			logx.L().Warnf("lease ctx done")
			_, err = lease.Revoke(context.TODO(), leaseID)
			if err != nil {
				logx.L().Errorf("lease revoke err:%v", err)
			}
			return
		}
	}
}
