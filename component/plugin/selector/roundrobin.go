// Author: huaxr
// Time:   2021/8/27 上午11:00
// Git:    huaxr

package selector

import "sync/atomic"

// roundRobinSelector selects servers with roundrobin.
type roundRobinSelector struct {
	servers []string
	r       *int32
}

func (s roundRobinSelector) Select() string {
	if len(s.servers) == 0 {
		return ""
	}
	i := *s.r
	i = i % int32(len(s.servers))

	atomic.AddInt32(s.r, 1)
	if *s.r >= int32(len(s.servers)) {
		atomic.StoreInt32(s.r, 0)
	}
	return s.servers[i]
}

func newRoundRobinSelector(servers map[string]string) Selector {
	var ss = make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}
	r := int32(0)
	return &roundRobinSelector{servers: ss, r: &r}
}
