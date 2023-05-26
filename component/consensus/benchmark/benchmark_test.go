// Author: huaxr
// Time:   2021/9/8 下午4:58
// Git:    huaxr

package benchmark

import (
	"context"
	"github.com/huaxr/magicflow/component/consensus"
	"testing"
	"time"
)

func BenchmarkGenID(b *testing.B) {
	consensus.LaunchIdGenerate(context.Background(), consensus.ETCD)
	time.Sleep(3 * time.Second)
	b.Run("a", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id := consensus.GetID()
			b.Logf("get id: %v", id)
		}
	})
	b.StopTimer()
}
