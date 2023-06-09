// Author: huaxr
// Time:   2021/8/27 上午11:01
// Git:    huaxr

package selector

import (
	"fmt"
	"github.com/edwingeng/doublejump"
	"hash/fnv"
	"sort"
)

// consistentHashSelector selects based on JumpConsistentHash.
type consistentHashSelector struct {
	h       *doublejump.Hash
	servers []string
}

// HashString get a hash value of a string
func HashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func toString(obj interface{}) string {
	return fmt.Sprintf("%v", obj)
}

func genKey(options ...interface{}) uint64 {
	panic("this implement not support")
	keyString := ""
	for _, opt := range options {
		keyString = keyString + "/" + toString(opt)
	}

	return HashString(keyString)
}

func (s consistentHashSelector) Select() string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}

	key := genKey("consistent key")
	return s.h.Get(key).(string)
}

func newConsistentHashSelector(servers map[string]string) Selector {
	var ss = make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	sort.Slice(ss, func(i, j int) bool { return ss[i] < ss[j] })
	return &consistentHashSelector{servers: ss, h: doublejump.NewHash()}
}
