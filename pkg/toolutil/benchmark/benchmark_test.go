// Author: huaxr
// Time:   2021/9/8 下午4:58
// Git:    huaxr

package benchmark

import (
	"github.com/huaxr/magicflow/pkg/toolutil"
	"testing"
)

func BenchmarkString2Byte(b *testing.B) {
	elems := "abcdefg"

	b.Run("common", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = []byte(elems)
		}
	})

	b.Run("zerocopy", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = toolutil.String2Byte(elems)
		}

	})

	b.StopTimer()
}
