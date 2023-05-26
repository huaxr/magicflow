// Author: huaxr
// Time:   2021/12/16 下午4:28
// Git:    huaxr

package taskService

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"time"
)

func CreateTask(c *gin.Context, req *CreateNormalTaskReq) (interface{}, error) {
	dao := orm.GetEngine()

	// 校验 name 字段
	if err := toolutil.AZaz09_(req.TaskName); err != nil {
		return nil, errors.New(fmt.Sprintf("task_name validate fail: %v", err.Error()))
	}

	// 校验 Configuration 字段
	var config map[string]interface{}
	if len(req.Configuration) > 0 {
		err := json.Unmarshal(toolutil.String2Byte(req.Configuration), &config)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Configuration must be a map or null"))
		}
	}

	// appid 要存在
	var apps = make([]models.App, 0)
	dao.Slave().Where("id = ?", req.AppId).Find(&apps)
	if len(apps) == 0 {
		return nil, errors.New(fmt.Sprintf("appid %d not exist", req.AppId))
	}

	// app必须经过核审
	if apps[0].Checked != 1 {
		return nil, errors.New(fmt.Sprintf("app %d not checked yet", req.AppId))
	}

	// 不能有重重复的任务
	var tasks = make([]models.Tasks, 0)
	dao.Slave().Where("name = ? and app_id = ?", req.Name, req.AppId).Find(&tasks)
	if len(tasks) > 0 {
		return nil, errors.New(fmt.Sprintf("task name %v already exist", req.Name))
	}

	// 鉴权
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err := user.HasAppPermit(cast.ToInt(req.AppId), false)
	if err != nil {
		return nil, err
	}

	insert := models.Tasks{
		Configuration: req.Configuration,
		Name:          req.Name,
		Description:   req.Description,
		AppId:         req.AppId,
		Xrn:           core.NewNormalTaskRN(req.TaskName),
		Type:          "normal",
		UpdateTime:    time.Now(),
		InputExample:  req.InputExample,
		OutputExample: req.OutputExample,
		User:          user.Account,
	}
	_, err = dao.Master().Insert(&insert)
	if err != nil {
		return CreateTaskRes{}, err
	}
	return CreateTaskRes{insert.Id}, nil
}

func CreateLocalTask(c *gin.Context, req *CreateLocalTaskReq) (interface{}, error) {
	dao := orm.GetEngine()
	// appid 要存在
	var apps = make([]models.App, 0)
	dao.Slave().Where("id = ?", req.AppId).Find(&apps)
	if len(apps) == 0 {
		return nil, errors.New(fmt.Sprintf("appid %d not exist", req.AppId))
	}

	// app必须经过核审
	if apps[0].Checked != 1 {
		return nil, errors.New(fmt.Sprintf("app %d not checked yet", req.AppId))
	}

	// 不能有重重复的任务
	var tasks = make([]models.Tasks, 0)
	dao.Slave().Where("name = ? and app_id = ?", req.Name, req.AppId).Find(&tasks)
	if len(tasks) > 0 {
		return nil, errors.New(fmt.Sprintf("task name %v already exist", req.Name))
	}

	// 查看本地剧本id是否存在
	var pbs = make([]models.Playbook, 0)
	dao.Slave().Where("id = ? and app_id = ?", req.PlaybookId, req.AppId).Find(&pbs)
	if len(pbs) == 0 {
		return nil, errors.New(fmt.Sprintf("local playbook id %v:%v not exist", req.PlaybookId, req.AppId))
	}

	// 鉴权
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err := user.HasAppPermit(cast.ToInt(req.AppId), false)
	if err != nil {
		return nil, err
	}

	insert := models.Tasks{
		Configuration: req.Configuration,
		Name:          req.Name,
		Description:   req.Description,
		AppId:         req.AppId,
		Xrn:           core.NewLocalTaskRN(req.PlaybookId),
		Type:          "local",
		UpdateTime:    time.Now(),
		InputExample:  req.InputExample,
		OutputExample: req.OutputExample,
		User:          user.Account,
	}
	_, err = dao.Master().Insert(&insert)
	if err != nil {
		return CreateTaskRes{}, err
	}
	return CreateTaskRes{insert.Id}, nil
}

func CreateRemoteTask(c *gin.Context, req *CreateRemotePlaybookTaskReq) (interface{}, error) {
	dao := orm.GetEngine()
	// appid 要存在
	var apps = make([]models.App, 0)
	dao.Slave().Where("id = ?", req.AppId).Find(&apps)
	if len(apps) == 0 {
		return nil, errors.New(fmt.Sprintf("appid %d not exist", req.AppId))
	}

	// app必须经过核审
	if apps[0].Checked != 1 {
		return nil, errors.New(fmt.Sprintf("app %d not checked yet", req.AppId))
	}

	// 判断任务是否已经存在
	var tasks = make([]models.Tasks, 0)
	dao.Slave().Where("name = ? and app_id = ?", req.Name, req.AppId).Find(&tasks)
	if len(tasks) > 0 {
		return nil, errors.New(fmt.Sprintf("task name %v already exist", req.Name))
	}

	// 找到 remote app 的id
	var app2 = make([]models.App, 0)
	dao.Slave().Where("app_name = ?", req.RemoteApp).Find(&app2)
	if len(app2) != 1 {
		return nil, errors.New(fmt.Sprintf("appname %s not exist", req.RemoteApp))
	}

	// 查看改远程id是否存在
	var pbs = make([]models.Playbook, 0)
	dao.Slave().Where("id = ? and app_id = ?", req.PlaybookId, app2[0].Id).Find(&pbs)
	if len(pbs) == 0 {
		return nil, errors.New(fmt.Sprintf("remote playbook %v:%v not exist", req.PlaybookId, req.RemoteApp))
	}

	// 鉴权
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err := user.HasAppPermit(cast.ToInt(req.AppId), false)
	if err != nil {
		return nil, err
	}

	insert := models.Tasks{
		Configuration: req.Configuration,
		Name:          req.Name,
		Description:   req.Description,
		AppId:         req.AppId,
		Xrn:           core.NewRemoteTaskRN(req.RemoteApp, req.PlaybookId),
		Type:          "remote",
		UpdateTime:    time.Now(),
		InputExample:  req.InputExample,
		OutputExample: req.OutputExample,
		User:          user.Account,
	}

	_, err = dao.Master().Insert(&insert)
	if err != nil {
		return CreateTaskRes{}, err
	}
	return CreateTaskRes{insert.Id}, nil
}
