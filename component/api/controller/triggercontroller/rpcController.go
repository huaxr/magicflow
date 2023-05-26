// Author: XinRui Hua
// Time:   2022/3/21 下午8:12
// Git:    huaxr

package triggercontroller

import (
	"encoding/json"
	"fmt"

	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/core"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	"github.com/huaxr/magicflow/component/api/controller"
	"github.com/huaxr/magicflow/component/service/client"
	"github.com/huaxr/magicflow/component/service/proto"
	"github.com/huaxr/magicflow/pkg/request"
)

type RpcController struct {
	controller.BaseApiController
}

func (ctl *RpcController) Router(router *gin.Engine) {
	entry := router.Group("/trigger")
	entry.POST("/execute", ctl.TriggerCall)
	entry.POST("/worker_response", ctl.WorkerResponseCall)
	entry.POST("/worker_exception", ctl.WorkerExceptionCall)
	entry.POST("/hook", ctl.HookStateMachineCall)
}

func (ctl *RpcController) HookStateMachineCall(c *gin.Context) {
	req := &request.HookStatePlaybook{}
	if err := c.ShouldBind(req); err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	// firstly, we should get mod parameters.
	dao := orm.GetEngine()
	var ns = make([]models.Execution, 0)
	table := fmt.Sprintf("execution_%d", cast.ToInt(req.SnapshotId)%10)
	err := dao.Slave().Table(table).Where("trace_id = ? and node_code = ? and status = ?", req.TraceId, req.NodeCode, core.Hooked).Find(&ns)
	if err != nil {
		logx.L().Errorf("HookPlaybook err %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	if len(ns) == 0 {
		ctl.Error(c, 200, "no result", nil)
		return
	}

	var extra core.Extra
	_ = json.Unmarshal([]byte(ns[0].Extra), &extra)

	srv, mod := extra.Detail.Meta.ServiceAddr, extra.Detail.Meta.Mod
	conn, err := client.GetConn(srv)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	client := proto.NewTriggerRpcClient(conn)

	_, err = client.Hook(c, &proto.HookReq{
		Srv:        srv,
		Mod:        mod,
		TraceId:    uint64(req.TraceId),
		NodeCode:   req.NodeCode,
		SnapshotId: uint32(req.SnapshotId),
	})
	if err != nil {
		logx.L().Errorf("HookStateMachineCall err %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	ctl.Success(c, gin.H{"data": "HookStateMachineCall success"})
}

func (ctl *RpcController) WorkerExceptionCall(c *gin.Context) {
	req := &request.WorkerExceptionReq{}
	if err := c.ShouldBind(req); err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	srv, mod := req.ServiceAddr, req.Key.GetMod()
	conn, err := client.GetConn(srv)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	b, _ := json.Marshal(req)
	client := proto.NewTriggerRpcClient(conn)
	_, err = client.WorkerException(c, &proto.WorkerExceptionReq{
		Mod:     uint32(mod),
		Srv:     srv,
		Payload: b,
	})
	if err != nil {
		logx.L().Errorf("WorkerExceptionCall err %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	ctl.Success(c, gin.H{"data": "WorkerExceptionCall success"})
}

func (ctl *RpcController) WorkerResponseCall(c *gin.Context) {
	req := &request.WorkerResponseReq{}
	if err := c.ShouldBind(req); err != nil {
		logx.L().Errorf("WorkerResponseCall bind WorkerResponseReq err, %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	srv, mod := req.ServiceAddr, req.Key.GetMod()
	conn, err := client.GetConn(srv)
	if err != nil {
		logx.L().Errorf("WorkerResponseCall get connection err, %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	b, err := json.Marshal(req)
	if err != nil {
		logx.L().Errorf("WorkerResponseCall marshal request err, %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	cli := proto.NewTriggerRpcClient(conn)
	_, err = cli.WorkerResponse(c, &proto.WorkerRespReq{
		Mod:     uint32(mod),
		Srv:     srv,
		Payload: b,
	})
	if err != nil {
		logx.L().Errorf("WorkerResponseCall err %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	ctl.Success(c, gin.H{"data": "WorkerResponseCall success"})
}

func (ctl *RpcController) TriggerCall(c *gin.Context) {
	params := &request.TriggerPlaybook{}
	if err := c.ShouldBind(params); err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	srv, mod := client.GetSrvAndMod()
	// 随机选取
	conn, err := client.GetConn(srv)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	b, _ := json.Marshal(params.Data)

	cli := proto.NewTriggerRpcClient(conn)
	resp, err := cli.Trigger(c, &proto.TriggerReq{
		Srv:        srv,
		Mod:        mod,
		AppName:    params.AppName,
		PlaybookId: int32(params.PlaybookId),
		AppToken:   params.AppToken,
		Sync:       params.Sync,
		Data:       b,
	})
	if err != nil {
		logx.L().Errorf("TriggerCall rpc call err: %v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	var res = gin.H{"trace_id": resp.Data}
	if len(resp.SyncResult) != 0 {
		var result interface{}
		json.Unmarshal(resp.SyncResult, &result)
		res["sync_result"] = result
	}
	ctl.Success(c, res)
}
