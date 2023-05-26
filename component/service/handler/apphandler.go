// Author: XinRui Hua
// Time:   2022/3/24 上午11:22
// Git:    huaxr

package handler

import (
	"context"
	"errors"

	"github.com/huaxr/magicflow/component/api/service/appservice"
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/dcc"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/service/proto"
)

type AppService struct{}

func (t *AppService) UpdateAppInternal(ctx context.Context, req *proto.UpdateAppInternalReq) (*proto.UpdateAppInternalResponse, error) {
	dao := orm.GetEngine()

	if appservice.HandlerStatus(req.Status) == appservice.Reject {
		updateReject := models.App{
			Checked: -1,
		}
		dao.Slave().MustCols("checked").Where("id = ?", req.AppId).Update(&updateReject)
		return &proto.UpdateAppInternalResponse{AppId: req.AppId}, nil
	}

	if appservice.HandlerStatus(req.Status) != appservice.Accept {
		return &proto.UpdateAppInternalResponse{AppId: req.AppId}, errors.New("not allowed Status")
	}

	// todo: 当前eps是不是分布式限流，实际落到k8s eps需要用 req.Eps 除以 pod 总数量,需要动态获取pod数
	if req.Eps < 100 {
		req.Eps = 100
	}

	if len(req.Brokers) == 0 {
		return &proto.UpdateAppInternalResponse{AppId: req.AppId}, errors.New("brokers not allocation yet")
	}

	update := models.App{
		Brokers: req.Brokers,
		Eps:     int(req.Eps),
		Checked: 1,
	}
	_, err := dao.Master().Where("id = ?", req.AppId).Update(&update)
	if err != nil {
		logx.L().Errorf("update app is err: %+v", err)
		return &proto.UpdateAppInternalResponse{AppId: req.AppId}, err
	}

	// dcc update
	d := dcc.GetDcc()
	err = d.DccPutApp(int(req.AppId))
	if err != nil {
		logx.L().Errorf("DccPutApp err: %+v", err)
		return &proto.UpdateAppInternalResponse{AppId: req.AppId}, err
	}

	return &proto.UpdateAppInternalResponse{AppId: req.AppId}, nil
}
