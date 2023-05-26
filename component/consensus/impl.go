// Author: XinRui Hua
// Time:   2022/2/10 上午10:33
// Git:    huaxr

package consensus

import (
	"context"

	"github.com/coreos/etcd/clientv3"
)

type Consensus interface {
	Close()
	Get(key string) [][]byte
	Put(key, val string) (err error)
	Elect(ctx context.Context)

	// need update
	WatchKey(key string) clientv3.WatchChan
	Delete(key string) (resp *clientv3.DeleteResponse, err error)
	GetClient() *clientv3.Client
}
