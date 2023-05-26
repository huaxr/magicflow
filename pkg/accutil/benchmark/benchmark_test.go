// Author: huaxr
// Time:   2021/9/8 下午4:58
// Git:    huaxr

package benchmark

import (
	"github.com/huaxr/magicflow/pkg/accutil"
	"testing"
)

func BenchmarkContains(b *testing.B) {
	elems := []string{"a", "a", "b", "ccxxxxx      xxxxx", "d", "e", "f", "g", "h", "ccxxxxx  z    xxxxx", "bb", "cc", "mmm"}

	var look = "b"
	b.Run("plan9", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ok := accutil.ContainsStr(elems, look)
			if !ok {
				b.Log("err plan9")
			}
		}

	})

	b.Run("common", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ok := accutil.ContainsStr(elems, look)
			if !ok {
				b.Log("err common")
			}
		}
	})

	b.StopTimer()
}
