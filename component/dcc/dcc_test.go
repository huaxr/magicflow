// Author: huaxr
// Time:   2021/12/13 上午11:05
// Git:    huaxr

package dcc

import (
	"context"
	"fmt"
	"github.com/huaxr/magicflow/component/consensus"
	"github.com/huaxr/magicflow/component/kv"
	"github.com/spf13/cast"
	"testing"
	"time"
)

func Test(t *testing.T) {
	ctx := context.Background()
	LaunchDcc(ctx, consensus.ETCD)
	kv.LunchRedis(ctx)

	globalDcc.KVSet("key", "aaaaaaa", time.Second*5)
	v, ok := globalDcc.KVGet("key")
	t.Logf("%v %v", v, ok)

	select {}
}

func BenchmarkDcc(b *testing.B) {
	ctx := context.Background()
	LaunchDcc(ctx, consensus.ETCD)
	kv.LunchRedis(ctx)

	b.Run("BenchmarkContains", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i <= 100000; i++ {
			//for i := 0; i < b.N; i++ {
			globalDcc.KVSet(fmt.Sprintf("%v_%v", "keysxx", i), "1", time.Second*1)
			//}
		}
	})
	b.StopTimer()
}

func TestCast(t *testing.T) {
	x := 0
	z := cast.ToBool(x)
	t.Log(z)
}
