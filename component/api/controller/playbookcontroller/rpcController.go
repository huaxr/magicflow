// Author: XinRui Hua
// Time:   2022/3/24 上午9:50
// Git:    huaxr

package playbookcontroller

import (
	"encoding/json"
	"fmt"

	"github.com/huaxr/magicflow/component/helper/transport"

	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/api/controller"
	"github.com/huaxr/magicflow/component/api/service/playbookservice"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/service/client"
	"github.com/huaxr/magicflow/component/service/proto"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type RpcController struct {
	controller.BaseApiController
}

// rpc 转化的接口
func (ctl *RpcController) Router(router *gin.Engine) {
	entry := router.Group("/playbook")
	// 获取剧本后端数据结构详情
	entry.GET("/get_real_playbook", ctl.GetRealPlayBook)
	// 切换快照
	entry.POST("/switch_snapshot", ctl.SwitchSnapshot)
	// 提交剧本并保存 raw 和 real struct，并通知dcc
	entry.POST("/submit", ctl.SubmitPlaybook)
}

func (ctl *RpcController) SubmitPlaybook(c *gin.Context) {
	var reqjson map[string]interface{}
	err := c.ShouldBind(&reqjson)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	conn, err := client.GetRandomConn()
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	b, _ := json.Marshal(reqjson)
	client := proto.NewPlaybookRpcClient(conn)
	resp, err := client.SubmitPlayBook(transport.CtxToGRpcCtx(c), &proto.SubmitPlayBookReq{
		Body: b,
	})
	if err != nil {
		logx.L().Errorf("SubmitPlaybook err %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	ctl.Success(c, resp)
	return
}

func (ctl *RpcController) SwitchSnapshot(c *gin.Context) {
	req := &playbookservice.SwitchVersionReq{}
	if err := c.ShouldBind(req); err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	dao := orm.GetEngine()

	var res = make([]models.Snapshot, 0)
	dao.Slave().Where("id = ? and playbook_id = ?", req.SnapShotId, req.PlayBookId).Find(&res)

	if len(res) != 1 {
		ctl.Error(c, 200, fmt.Errorf("snapshot %v not exist", req.SnapShotId).Error(), nil)
		return
	}
	// 鉴权
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err := user.HasAppPermit(cast.ToInt(res[0].AppId), true)
	if err != nil {
		logx.L().Errorf("SwitchSnapshot err %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	conn, err := client.GetRandomConn()
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	client := proto.NewPlaybookRpcClient(conn)
	resp, err := client.SwitchSnapshot(c, &proto.SwitchVersionReq{
		PlayBookId: int32(req.PlayBookId),
		SnapShotId: int32(req.SnapShotId),
	})
	if err != nil {
		logx.L().Errorf("SwitchSnapshot err %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	ctl.Success(c, playbookservice.SwitchVersionResp{PlayBookId: int(resp.PlayBookId)})
}

func (ctl *RpcController) GetRealPlayBook(c *gin.Context) {
	smId, exist := c.GetQuery("id")
	if !exist || cast.ToInt(smId) <= 0 {
		logx.L().Errorf("id: %v not right!", smId)
		ctl.Error(c, 200, "err id", nil)
		return
	}
	req := playbookservice.GetPlayBookReq{
		PlayBookId: smId,
	}
	pb, err := playbookservice.GetRealPlayBook(c, &req)
	if err != nil {
		logx.L().Errorf("GetRealPlayBook is err: %+v!", err)
		ctl.Error(c, 200, "GetRealPlayBook is err", err)
		return
	}

	conn, err := client.GetRandomConn()
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	client := proto.NewPlaybookRpcClient(conn)
	resp, err := client.GetPlaybook(c, &proto.GetPlaybookReq{
		Pbid: int32(cast.ToInt(smId)),
	})

	if err != nil {
		logx.L().Errorf("HookStateMachineCall %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	var response map[string]interface{}
	json.Unmarshal(resp.Data, &response)
	response["description"] = pb.Description
	response["name"] = pb.Name
	response["snapshot_id"] = pb.SnapshotId
	ctl.Success(c, response)
}
