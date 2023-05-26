// Author: XinRui Hua
// Time:   2022/4/14 下午5:40
// Git:    huaxr

package promethu

import (
	"context"
	"testing"
)

func TestProm(t *testing.T) {
	LaunchProm(context.TODO(), Push)
}
