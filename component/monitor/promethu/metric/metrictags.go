// Author: XinRui Hua
// Time:   2022/4/11 下午6:04
// Git:    huaxr

package metric

import (
	"sync/atomic"

	"github.com/huaxr/magicflow/component/monitor/promethu/tag"
	"github.com/huaxr/magicflow/component/monitor/promethu/vars"
)

var (
	tags = make([]*m, 0)
)

type m struct {
	key   tag.TagKey
	count int32
}

func (x m) GetCount() int32 {
	return atomic.LoadInt32(&x.count)
}

func (x m) GetKey() string {
	return x.key.String()
}

func (x m) GetTag() tag.TagKey {
	return x.key
}

func RegisterMetrics() {
	key := []tag.TagKey{
		tag.Trigger,
		tag.WorkerException,
		tag.ServerException,
		tag.Complied,
		tag.PlaybookPut,
		tag.AppPut,
		tag.AppSwitch,
	}

	for index, i := range key {
		if index != int(i) {
			panic("registerMetrics sequence err")
		}
		tags = append(tags, &m{
			key:   i,
			count: 0,
		})
	}
}

// counter
func Metric(keys ...tag.TagKey) {
	if len(keys) == 0 {
		return
	}
	key := keys[0]

	switch key {
	case tag.Trigger:
		vars.IncResidualTask()
	case tag.Complied, tag.WorkerException, tag.ServerException:
		vars.DecResidualTask()
	}

	atomic.AddInt32(&tags[key].count, 1)
}

func Clean(tag tag.TagKey) {
	atomic.SwapInt32(&tags[tag].count, 0)
}

func GetMetric() []*m {
	cpy := make([]*m, len(tags))
	copy(cpy, tags)
	return cpy
}
