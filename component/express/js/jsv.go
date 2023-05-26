// Author: huaxr
// Time:   2021/6/7 上午11:22
// Git:    huaxr

package js

import (
	"encoding/json"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/robertkrimen/otto"
	"sync"
)

const (
	VmCount = 10
	Key     = "Context"
)

type Jsv struct {
	lock *sync.Mutex
	vm   *otto.Otto
}

var (
	jsvMap map[int]*Jsv
	once   sync.Once
)

type JsvExpress struct {
	Express   string
	Input     interface{}
	Condition string
}

func (j *JsvExpress) EExecute() (bool, error) {
	var res otto.Value
	var err error

	res, err = execute(j.Input, j.Condition)

	if err != nil {
		return false, err
	}
	b, err := res.ToBoolean()
	if err != nil {
		return false, err
	}
	return b, nil
}

func Init() {
	once.Do(func() {
		jsvMap = make(map[int]*Jsv)
		for i := 0; i <= VmCount-1; i++ {
			jsv := new(Jsv)
			jsv.vm = otto.New()
			jsv.lock = new(sync.Mutex)
			jsvMap[i] = jsv
		}
	})
}

func execute(input interface{}, condition string) (otto.Value, error) {
	if jsvMap == nil {
		Init()
	}
	b, _ := json.Marshal(input)
	index := toolutil.CRC(b) % VmCount
	jsvMap[index].lock.Lock()
	defer jsvMap[index].lock.Unlock()
	jsvMap[index].vm.Set(Key, input)
	return jsvMap[index].vm.Run(condition)
}

func (j *JsvExpress) GetExpress() string {
	return j.Express
}
