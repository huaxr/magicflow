// Author: XinRui Hua
// Time:   2022/5/6 上午10:16
// Git:    huaxr

package core

import (
	"testing"
)

func TestMessage(t *testing.T) {
	newTestPlaybook(t)

	m := message{
		Meta: meta{
			Trace:       123456,
			Mod:         1,
			ServiceAddr: "111",
			Topic:       "tmp",
			MessageType: ToServer,
			Sequence:    1,
			Timestamp:   "",
			Domain:      "",
			Signature:   nil,
			Sync:        false,
		},
		Task: task{
			Status:         Flying,
			PlaybookId:     1,
			SnapshotId:     8,
			NodeCode:       "",
			Input:          "1",
			Output:         "2",
			InputAttribute: nil,
			Xrn: &xrn{
				Service:   "",
				Region:    "",
				TaskType:  "",
				Namespace: "",
				TaskInfo:  "",
			},
			Configuration: nil,
		},
		Context: context{
			Last:      []interface{}{map[string]interface{}{"xxx": "yyy"}, 2, 3},
			Store:     map[string]interface{}{"a": []interface{}{"1", "a", "x"}, "x": "y"},
			Exception: exception{},
			Stack:     nil,
			Env:       map[string]interface{}{"a": "v"},
			Chain:     nil,
			Heartbeat: nil,
		},
	}
	res, err := m.attrHandler(map[string]interface{}{"key": "$.x", "key2": []interface{}{1, 2, "$$.[0].xxx"}, "key3": map[string]interface{}{"a": "$.a[0]"}})
	t.Log(res, err)

	res, err = m.attrHandler([]interface{}{"$$.[0].xxx", "$$.[1]", 2})
	t.Log(res, err)

	err = m.ResolveAmbiguous()
	t.Log(err)

}

type z struct {
	age []string
	m   map[string]int
}

func (z *z) push(a string) {
	z.age = append(z.age, a)
}

func (z z) add(a string, v int) {
	z.m[a] = v
}

type x struct {
	name string
	nums []int
	z    *z
	zz   z
}

func (x *x) copy() x {
	y := *x
	return y
}

func TestDeepCopy(t *testing.T) {
	var m = &x{"xxx", []int{1, 2, 3}, &z{[]string{}, map[string]int{}}, z{[]string{}, map[string]int{}}}
	y := m.copy()
	y.name = "yyy"
	y.nums = []int{4, 5, 6}

	m.z.push("c")
	m.z.add("aa", 1)
	y.z.push("d")
	y.z.add("bb", 1)

	m.zz.push("c")
	m.zz.add("cc", 1)
	y.zz.push("d")
	y.zz.add("dd", 1)

	t.Log(m, y)
	t.Log(m.z, y.z)
	t.Log(m.zz, y.zz)
}

func TestM(t *testing.T) {
	a := map[string]int{"A": 1}

	b := a
	b["v"] = 2

	t.Log(a, b)
}

func TestDele(t *testing.T) {
	m := map[string]*z{
		"z": &z{
			age: []string{"A", "B"},
			m:   nil,
		},
	}

	z := *m["z"]

	delete(m, "z")
	t.Log(&z.age)
}

func BenchmarkEccSign(b *testing.B) {
	msg := Msg{
		Key:         "abcdEFGHIGKLMNOPQRSTUvw",
		ServiceAddr: "xyz",
	}
	n := b.N
	b.Run("ECC", func(b *testing.B) {
		for i := 0; i < n; i++ {
			msg.Sign(Ecc)
			err := msg.CheckSign()
			if err != nil {
				b.Log("ecc err:", err)
			}
		}
	})

	b.Run("JWT", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < n; i++ {
			msg.Sign(Jwt)
			err := msg.CheckSign()
			if err != nil {
				b.Log("jwt err:", err)
			}
		}
	})

	b.StopTimer()
}
