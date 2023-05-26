package eccutil

import (
	"testing"

	"github.com/huaxr/magicflow/pkg/toolutil"
)

func TestEcc(t *testing.T) {

	b, _ := ECCEncrypt(toolutil.String2Byte("abcd"))

	res, _ := ECCDecrypt(b)
	t.Log(string(res))
}

func TestEccSign(t *testing.T) {
	pt := toolutil.String2Byte("abcdEFGHIGKLMNOPQRSTUvwxysssssssssssssssssssssssssssssssssssz")
	sig, _ := EccSign(pt)
	ok := EccSignVer(pt, sig)
	t.Log(ok)

	ok = EccSignVer(pt, append(sig, byte(1)))
	t.Log(ok)
}

func BenchmarkEccSign(b *testing.B) {
	pt := toolutil.String2Byte("abcdEFGHIGKLMNOPQRSTUvwxyz")

	b.Run("ECC", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sig, _ := EccSign(pt)
			EccSignVer(pt, sig)
		}
	})

	b.StopTimer()
}
