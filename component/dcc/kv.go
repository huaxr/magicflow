// Author: huaxr
// Time:   2021/12/30 下午7:22
// Git:    huaxr

package dcc

import (
	"encoding/json"
	"fmt"

	"github.com/huaxr/magicflow/component/logx"

	"reflect"
	"time"

	"github.com/huaxr/magicflow/component/kv"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/huaxr/magicflow/pkg/toolutil"
)

// Configuration notify
// group represents that
// implements KVCache in kv package
func (d *dcc) KVSet(groupKey string, value interface{}, duration time.Duration) {
	k := fmt.Sprintf("%s/%s/%s", confutil.GetConf().Configuration.WatchKeyPrefix, Configuration, groupKey)
	// update local & redis storage
	kv.RedisKv.KVSet(groupKey, value, duration)

	var putStr string
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		putStr = value.(string)
	default:
		b, _ := json.Marshal(value)
		putStr = toolutil.Bytes2string(b)
	}
	err := d.con.Put(k, putStr)
	if err != nil {
		logx.L().Errorf("put is err: %+v", err)
		return
	}
	logx.L().Infof("put success key:%+v,  value:%+v", groupKey, value)
}

func (d *dcc) KVGet(key string) (val interface{}, exist bool) {
	val, exist = kv.RedisKv.KVGet(key)
	if !exist {
		// or load from persistence storage like mysql.
		val = confutil.GetConf().Configuration.None
		if val == nil {
			return "", false
		}
		return val, true
	}
	return
}
