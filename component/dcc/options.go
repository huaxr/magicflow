// Author: huaxr
// Time:   2021/10/12 下午3:47
// Git:    huaxr

package dcc

import (
	"fmt"

	"github.com/huaxr/magicflow/component/logx"

	"github.com/huaxr/magicflow/pkg/confutil"
)

// The DccPutPb function is used to implement the put function of etcd
func (d *dcc) DccPutPb(pbId, ssId int) error {
	prefix := confutil.GetConf().Configuration.WatchKeyPrefix
	key := fmt.Sprintf("%s/%s/%d", prefix, PLAYBOOK, pbId)
	value := fmt.Sprintf("version_%d", ssId)

	err := d.con.Put(key, value)
	if err != nil {
		logx.L().Errorf("put is err: %+v", err)
		return err
	}
	logx.L().Infof("put playbook success key:%+v,  value:%+v", key, value)
	return nil
}

func (d *dcc) DccPutApp(appId int) error {
	prefix := confutil.GetConf().Configuration.WatchKeyPrefix
	key := fmt.Sprintf("%s/%s/%d", prefix, APP, appId)
	value := fmt.Sprintf("%d", appId)

	err := d.con.Put(key, value)
	if err != nil {
		logx.L().Errorf("DccPutApp put is err: %+v", err)
		return err
	}
	logx.L().Infof("DccPutApp put app success key:%+v,  value:%+v", key, value)
	return nil
}

func (d *dcc) DccPutTask(taskid int) error {
	key := fmt.Sprintf("%s/%s/%d", confutil.GetConf().Configuration.WatchKeyPrefix, Task, taskid)
	value := fmt.Sprintf("%d", taskid)

	err := d.con.Put(key, value)
	if err != nil {
		logx.L().Errorf("DccPutTask put is err: %+v", err)
		return err
	}
	logx.L().Infof("DccPutTask put task success key:%+v,  value:%+v", key, value)
	return nil
}

type ChannelStatus string

const (
	OPEN  ChannelStatus = "open"
	CLOSE ChannelStatus = "close"
)

func (d *dcc) DccSwitchApp(appName string, status ChannelStatus) error {
	key := fmt.Sprintf("%s/%s/%s", confutil.GetConf().Configuration.WatchKeyPrefix, SwitchApp, appName)
	value := status

	err := d.con.Put(key, string(value))
	if err != nil {
		logx.L().Errorf("DccSwitchApp put is err: %+v", err)
		return err
	}
	logx.L().Infof("DccSwitchApp put switch app success key:%+v,  value:%+v", key, value)
	return nil
}

// The DccDelPb function is used to implement the delete function of etcd
func (d *dcc) DccDelPb(pbId int) error {
	key := fmt.Sprintf("%s/%s/%d", confutil.GetConf().Configuration.WatchKeyPrefix, PLAYBOOK, pbId)
	re, err := d.con.Delete(key)
	if err != nil {
		logx.L().Errorf("DccDelPb Delete is err: %+v", err)
		return err
	}
	logx.L().Infof("DccDelPb Delete pb success key:%+v  %+v", key, re)
	return nil
}
