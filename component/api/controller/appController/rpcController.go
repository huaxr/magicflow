// Author: XinRui Hua
// Time:   2022/3/24 上午11:18
// Git:    huaxr

package appController

import (
	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/api/controller"
	"github.com/huaxr/magicflow/component/api/service/appservice"
	"github.com/huaxr/magicflow/component/service/client"
	"github.com/huaxr/magicflow/component/service/proto"
	"github.com/gin-gonic/gin"
)

type RpcController struct {
	controller.BaseApiController
}

func (ctl *RpcController) Router(router *gin.Engine) {
	entry := router.Group("/app")
	// 管理员审核通过app（触发dcc）
	entry.POST("/admin_audit", ctl.UpdateAppInternal)
}

func (ctl *RpcController) UpdateAppInternal(c *gin.Context) {
	var params appservice.UpdateAppInternalReq
	err := c.ShouldBind(&params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	uif, _ := c.Get("auth")
	if uif.(*auth.UserInfo).IsAdmin != 1 {
		ctl.Error(c, 200, "you are not super admin", nil)
		return
	}

	conn, err := client.GetRandomConn()
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	client := proto.NewAppRpcClient(conn)
	resp, err := client.UpdateAppInternal(c, &proto.UpdateAppInternalReq{
		Status:  string(params.Status),
		AppId:   int32(params.AppId),
		Brokers: params.Brokers,
		Eps:     int32(params.Eps),
	})
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	ctl.Success(c, resp)
	return
}
