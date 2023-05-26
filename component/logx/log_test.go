// Author: XinRui Hua
// Time:   2022/4/5 上午10:10
// Git:    huaxr

package logx

import (
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	L(zap.String("x", "Z"), zap.Int("xz", 1)).Infof("xxx")
}
