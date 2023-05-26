// Author: huaxr
// Time:   2021/10/12 下午3:48
// Git:    huaxr

package dcc

import (
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/core"
)

func updateApp(appid int) {
	app := getappbyid(appid)
	if app == nil {
		return
	}
	core.SetNamespace(app)

	return
}

func closeApp(appName string) {
	ns, ok := core.G().GetNamespace(appName)
	if !ok {
		logx.L().Errorf("closeApp %v appName not exist", appName)
		return
	}
	if ns.Working() {
		ns.Close()
	}
}

func openApp(appName string) {
	ns, ok := core.G().GetNamespace(appName)
	if !ok {
		logx.L().Errorf("openApp %v appName not exist", appName)
		return
	}
	if !ns.Working() {
		ns.Open()
	}
}

func getappbyid(appid int) *models.App {
	dao := orm.GetEngine()

	res := make([]models.App, 0)
	err := dao.Slave().Where("id = ?", appid).Find(&res)
	if err != nil {
		logx.L().Errorf("search app by id err:%+v", err)
		return nil
	}
	if len(res) != 1 {
		logx.L().Errorf("app id is not exist or not unique", err)
		return nil
	}

	return &res[0]
}

func updateTask(taskid int) {
	dao := orm.GetEngine()
	res := make([]models.Tasks, 0)
	err := dao.Slave().Where("id = ?", taskid).Find(&res)
	if err != nil {
		logx.L().Errorf("search Tasks by id err:%+v", err)
		return
	}
}
