// Author: huaxr
// Time:   2021/9/8 下午4:58
// Git:    huaxr

package benchmark

import (
	"github.com/huaxr/magicflow/component/express/parser"
	"testing"
)

func BenchmarkContains(b *testing.B) {
	b.Run("BenchmarkContains", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			parser.GetKey("$.AAA.BBB", "", map[string]interface{}{"AAA": map[string]interface{}{"BBB": "XXX"}})
		}
	})

	b.Run("BenchmarkContains", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			parser.GetKey("$.AAA[0]", "", map[string]interface{}{"AAA": []interface{}{0, 1, 2}})
		}
	})
	b.StopTimer()
}
