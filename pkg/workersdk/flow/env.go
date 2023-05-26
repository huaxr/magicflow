// Author: XinRui Hua
// Time:   2022/4/27 上午10:35
// Git:    huaxr

package flow

import "github.com/huaxr/magicflow/component/logx"

type key struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Next  *key        `json:"next"`
}

type chainenv struct {
	Keys key `json:"keys"`
	// recording the sequence of node chain.
	// provide retrospect function point
	// Chain []string
	Size int `json:"size"`
}

// insert at the front top
func (e *chainenv) Set(keys string, val interface{}) {
	var tmp key
	if e.Size >= 100 {
		logx.L().Warnf("env.set", "chain size limited, option ignored")
		return
	}
	tmp = e.Keys
	k := key{
		Key:   keys,
		Value: val,
		Next:  &tmp,
	}
	e.Keys = k
	e.Size += 1
}

func (e *chainenv) Get(key string) (val interface{}, exist bool) {
	tmp := e.Keys
	for {
		if tmp.Key == key {
			return tmp.Value, true
		} else {
			if tmp.Next == nil {
				break
			}
			tmp = *tmp.Next
		}
	}
	return nil, false
}
