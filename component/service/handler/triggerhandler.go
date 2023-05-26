// Author: XinRui Hua
// Time:   2022/3/18 下午4:38
// Git:    huaxr

package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/huaxr/magicflow/component/monitor/promethu/metric"
	"github.com/huaxr/magicflow/component/monitor/promethu/tag"

	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/spf13/cast"

	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/service/proto"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/request"
)

type TriggerService struct{}

func (t *TriggerService) Trigger(ctx context.Context, req *proto.TriggerReq) (*proto.TriggerResponse, error) {
	namespace, ok := core.G().GetNamespace(req.AppName)
	if !ok {
		return nil, fmt.Errorf("appname not exist")
	}
	if !namespace.Working() {
		return nil, fmt.Errorf("app %v shutdown currently cause there is no consumer found, wait for a second please", req.AppName)
	}

	pb, ok := core.G().GetPlaybook(int(req.PlaybookId))
	if !ok {
		return nil, fmt.Errorf("invalid sm id")
	}

	if pb.HasHook() && req.Sync {
		return nil, fmt.Errorf("playbook with hook could not use sync")
	}

	if pb.GetAppToken() != req.AppToken {
		return nil, fmt.Errorf("invalid token")
	}

	if pb.GetApp() != req.AppName {
		return nil, fmt.Errorf("invalid namespace with pbid:%v and namespace: %v", req.PlaybookId, req.AppName)
	}

	ns, ok := core.G().GetNamespace(req.AppName)
	if !ok {
		return nil, fmt.Errorf("app not exist")
	}
	limit := ns.GetLimiter()
	if limit != nil && !limit.Request() {
		return nil, fmt.Errorf("eps limitation is: %d", limit.GetQuota())
	}

	// 将byte再转化为 interface
	var data interface{}
	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		return &proto.TriggerResponse{
			Data: err.Error(),
		}, err
	}

	tri := core.NewTriggerMessage(int(req.PlaybookId), req.AppName, data, req.Srv, req.Mod, req.Sync)
	tri.Dispatch()
	metric.Metric(tag.Trigger)
	res, err := tri.WaitSync(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.TriggerResponse{
		Data:       cast.ToString(tri.Meta.Trace),
		SyncResult: res,
	}, nil
}

func (t *TriggerService) WorkerResponse(ctx context.Context, req *proto.WorkerRespReq) (*proto.WorkerRespResponse, error) {
	var param request.WorkerResponseReq
	err := json.Unmarshal(req.Payload, &param)
	if err != nil {
		return nil, err
	}
	msg, err := core.GetExchange().Ack(param.Key)
	if err != nil {
		logx.L().Errorf("WorkerResponse Ack err:%v", err)
		return nil, err
	}
	// 首先检验快照版本歧义
	err = msg.ResolveAmbiguous()
	if err != nil {
		logx.L().Errorf("WorkerResponse ResolveAmbiguous err:%v", err)
		return nil, err
	}

	tmp := core.Msg{
		Key:         param.Key,
		ServiceAddr: param.ServiceAddr,
		Signature:   param.Signature,
	}

	err = tmp.CheckSign()
	if err != nil {
		logx.L().Errorf("WorkerResponse check sign err:%v", err)
		return nil, err
	}

	srv := msg.NewServerMessage(param.Output, param.Env, param.HeartBeat)
	srv.Dispatch()

	return &proto.WorkerRespResponse{
		Data: "success",
	}, nil
}

func (t *TriggerService) WorkerException(ctx context.Context, req *proto.WorkerExceptionReq) (*proto.WorkerExceptionResponse, error) {
	var param request.WorkerExceptionReq
	json.Unmarshal(req.Payload, &param)

	msg, err := core.GetExchange().Ack(param.Key)
	if err != nil {
		return nil, err
	}
	//首先检验快照版本歧义
	err = msg.ResolveAmbiguous()
	if err != nil {
		return nil, err
	}

	exp := msg.NewExceptionMessage(param.Exception)
	exp.Dispatch()
	return &proto.WorkerExceptionResponse{
		Data: "WorkerException success",
	}, nil
}

func (t *TriggerService) Hook(ctx context.Context, req *proto.HookReq) (*proto.HookResponse, error) {
	traceId := req.TraceId
	nodeCode := req.NodeCode

	dao := orm.GetEngine()

	var ns = make([]models.Execution, 0)
	table := fmt.Sprintf("execution_%d", cast.ToInt(req.SnapshotId)%10)
	err := dao.Slave().Table(table).Where("trace_id = ? and node_code = ? and status = ?", traceId, nodeCode, core.Hooked).Find(&ns)
	if err != nil {
		logx.L().Errorf("HookPlaybook err %v", err)
		return nil, err
	}

	if len(ns) == 0 {
		return nil, fmt.Errorf("query execution not exist")
	}

	var extra core.Extra
	_ = json.Unmarshal([]byte(ns[0].Extra), &extra)

	// 更新hook的状态为 executed
	go extra.Detail.UpdateStatus(core.Executed)

	// 检查版本歧义
	err = extra.Detail.ResolveAmbiguous()
	if err != nil {
		logx.L().Errorf("ResolveAmbiguous for hook err, %v", err)
		return nil, err
	}
	// not need check hook, just go through.
	hatches := extra.Detail.Hatch(false)
	if len(hatches) > 0 {
		for _, i := range hatches {
			x := i
			x.Dispatch()
		}
	}

	return &proto.HookResponse{
		Data: "Hook Success",
	}, nil
}
