/**
 * @Author: huaxr
 * @Description:
 * @File: appService
 * @Version: 1.0.0
 * @Date: 2021/9/7 下午3:41
 */

package appservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/gin-gonic/gin"
	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/dcc"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"time"
)

func CreateApp(ctx *gin.Context, req *CreateAppReq) (CreateAppResp, error) {
	dao := orm.GetEngine()

	if len(req.AppName) < 5 || len(req.AppName) > 30 {
		return CreateAppResp{}, errors.New("appname length must between 5 and 30")
	}

	if err := toolutil.AZaz09_(req.AppName); err != nil {
		return CreateAppResp{}, errors.New(fmt.Sprintf("app_name validate fail: %v", err.Error()))
	}

	// Check the uniqueness of the app_name
	var res = make([]models.App, 0)
	err := dao.Slave().Where("app_name = ?", req.AppName).Find(&res)
	if err != nil {
		logx.L().Errorf("CreateApp.err", "find app by app_name is err: %+v", err)
		return CreateAppResp{}, err
	}
	if len(res) != 0 {
		logx.L().Errorf("CreateApp.err", "app_name already exists")
		return CreateAppResp{}, errors.New("the name already exists")
	}

	uif, _ := ctx.Get("auth")
	name := uif.(*auth.UserInfo).Account

	// Create new app
	create := models.App{
		AppName: req.AppName,
		Token:   toolutil.Base64Encode(toolutil.String2Byte(toolutil.GetRandomString(15))),
		Eps:     req.Eps,
		User:    name,
		GroupId: "",
		// default brokers
		Brokers:       "",
		UpdateTime:    time.Now(),
		Checked:       0,
		Share:         0,
		Description:   req.Description,
		LastAliveTime: time.Now().AddDate(0, 0, -2),
	}

	_, err = dao.Master().Insert(&create)
	if err != nil {
		logx.L().Errorf("CreateApp.err", "creat app is err: %+v", err)
		return CreateAppResp{}, err
	}
	id := create.Id

	return CreateAppResp{id}, nil
}

func UpdateAppInternal(ctx context.Context, req *UpdateAppInternalReq) (UpdateAppInternalResp, error) {
	dao := orm.GetEngine()
	if req.Status == Reject {
		updateReject := models.App{
			Checked: -1,
		}
		dao.Master().MustCols("checked").Where("id = ?", req.AppId).Update(&updateReject)
		return UpdateAppInternalResp{req.AppId}, nil
	}

	if req.Status != Accept {
		return UpdateAppInternalResp{req.AppId}, errors.New("not allowed Status")
	}

	// todo: 当前eps是不是分布式限流，实际落到k8s eps需要用 req.Eps 除以 pod 总数量,需要动态获取pod数
	if req.Eps < 100 {
		req.Eps = 100
	}

	if len(req.Brokers) == 0 {
		return UpdateAppInternalResp{req.AppId}, errors.New("brokers not allocation yet")
	}

	update := models.App{
		Brokers: req.Brokers,
		Eps:     req.Eps,
		Checked: 1,
	}
	_, err := dao.Master().Where("id = ?", req.AppId).Update(&update)
	if err != nil {
		logx.L().Errorf("UpdateAppInternal.err", "update app is err: %+v", err)
		return UpdateAppInternalResp{req.AppId}, err
	}

	// dcc update
	d := dcc.GetDcc()
	err = d.DccPutApp(req.AppId)
	if err != nil {
		logx.L().Errorf("appservice.UpdateAppInternal.err", "DccPutApp err: %+v", err)
		return UpdateAppInternalResp{req.AppId}, err
	}

	return UpdateAppInternalResp{req.AppId}, nil
}
