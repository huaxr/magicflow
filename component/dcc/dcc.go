// Author: huaxr
// Time:   2021/8/2 下午7:09
// Git:    huaxr

package dcc

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/huaxr/magicflow/pkg/confutil"

	"github.com/huaxr/magicflow/component/monitor/promethu/tag"

	"github.com/huaxr/magicflow/component/monitor/promethu/metric"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/huaxr/magicflow/component/consensus"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/core"
	"github.com/spf13/cast"
)

type (
	dcc struct {
		ctx  context.Context
		con  consensus.Consensus
		once sync.Once
	}

	WATCH string
)

var globalDcc *dcc

const (
	PLAYBOOK WATCH = "playbook"
	APP      WATCH = "app"
	// task update will regenerate all the playbook will
	// reference it.
	Task WATCH = "task"
	// playbooks in the phcache only need update or delete.
	// but app provides a monitor for it's consumers' alive
	// status, which need update app's enable field.
	SwitchApp     WATCH = "switch_app"
	Configuration WATCH = "configration" //  configration/moudle/key
)

func (w WATCH) String() string {
	return string(w)
}

// The dccWatch function is used to implement the watch function of etcd
// Check the occurrence of events and change the memory of each machine in time
func (d *dcc) dccWatch() {
	watchChan := d.con.WatchKey(confutil.GetConf().Configuration.WatchKeyPrefix)

WATCH:
	for i := range watchChan {
		// metric here
		for _, e := range i.Events {
			k := string(e.Kv.Key)
			key := strings.Split(k, "/")
			if len(key) != 5 {
				logx.L().Errorf("etcd key format error, value:%+v", key)
				continue WATCH
			}

			switch WATCH(key[3]) {
			case Configuration:

			case PLAYBOOK:
				// /magicFlow/watch/playbook/1   smid_snapid
				smId := cast.ToInt(key[4])
				logx.L().Infof("switch_pb id:%v, val:%v", smId, ChannelStatus(e.Kv.Value))
				switch e.Type {
				case mvccpb.PUT:
					metric.Metric(tag.PlaybookPut)
					v := string(e.Kv.Value)
					value := strings.Split(v, "_")
					if len(value) != 2 {
						logx.L().Errorf("etcd value format error, value:%+v", value)
						continue WATCH
					}
					snapshotId, err := strconv.Atoi(value[1])
					if err != nil {
						logx.L().Errorf("snapshotId strconv.Atoi is err: %+v", err)
						continue WATCH
					}
					err = core.ReloadSnapshot(snapshotId, smId)
					if err != nil {
						logx.L().Errorf("dccWatch err %v", err)
						continue WATCH
					}
				}

			case APP:
				// /magicFlow/watch/app/1   1
				appId := cast.ToInt(key[4])
				logx.L().Infof("switch_app watch id:%v, val:%v", appId, ChannelStatus(e.Kv.Value))

				switch e.Type {
				case mvccpb.PUT:
					metric.Metric(tag.AppPut)

					updateApp(appId)
				}

			case SwitchApp:
				appName := key[4]
				logx.L().Infof("switch_app watch %v %v", appName, ChannelStatus(e.Kv.Value))
				switch e.Type {
				case mvccpb.PUT:
					metric.Metric(tag.AppSwitch)
					switch ChannelStatus(e.Kv.Value) {
					case OPEN:
						openApp(appName)
					case CLOSE:
						closeApp(appName)
					default:

					}
				}
			}
		}
	}
}
