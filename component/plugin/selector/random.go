// Author: huaxr
// Time:   2021/8/27 上午10:59
// Git:    huaxr

package selector

import "github.com/valyala/fastrand"

// randomSelector selects randomly.
type randomSelector struct {
	servers []string
}

func (s randomSelector) Select() string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}
	i := fastrand.Uint32n(uint32(len(ss)))
	return ss[i]
}

func newRandomSelector(servers map[string]string) Selector {
	var ss = make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}
	return &randomSelector{servers: ss}
}
