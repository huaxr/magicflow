// Author: huaxr
// Time:   2021/12/21 下午2:31
// Git:    huaxr

package healthy

import (
	"testing"
)

func TestReportChannelStats(t *testing.T) {
	x := ReportChannelStats("my_test_worker")
	t.Logf("%v", x)
}
