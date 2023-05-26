// Author: huaxr
// Time:   2021/9/27 下午4:50
// Git:    huaxr

package healthy

import (
	"context"
)

type HealthyImpl interface {
	Report()
	String() string
}

var healthy = [...]struct {
	stat HealthyImpl
}{
	//{stat: &Profile{}},
	{stat: &channelStats{}},
}

func LaunchHealthy(ctx context.Context) {
	for _, i := range healthy {
		task := i
		go task.stat.Report()
	}
}
