// Author: XinRui Hua
// Time:   2022/5/5 下午5:22
// Git:    huaxr

package core

import "testing"

func TestTTL(t *testing.T) {
	// {"key":"$.bf9048.list[0]", "key2":[1,2,3], "key3":{"x":"y"}}
	v := map[string]interface{}{"key": "$.bf9048.list[0]", "key2": []interface{}{1, " $.xxxxxx.a", "$$.aaa"}, "key3": map[string]interface{}{"a": "$.aaaaaa.list[1]"}, "key4": []interface{}{1, "2", 3}}
	var tt = make(ttl)
	r := ttlx(v, tt)
	t.Log(r)
}
