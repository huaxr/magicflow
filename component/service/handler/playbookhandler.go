// Author: XinRui Hua
// Time:   2022/3/24 上午10:13
// Git:    huaxr

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/huaxr/magicflow/component/logx"

	"github.com/huaxr/magicflow/component/plugin/limiter"

	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/dcc"
	"github.com/huaxr/magicflow/component/service/proto"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

type PlaybookService struct {
	Name  string
	Limit limiter.Limiter
}

// CheckSum function checks whether body has been loaded
// true means body has been loaded, we should refuse
func CheckSumExist(sum string) (bool, int) {
	dao := orm.GetEngine()

	var res = make([]models.Snapshot, 0)
	err := dao.Slave().Where("checksum = ?", sum).Find(&res)
	if err != nil {
		logx.L().Errorf("find err:%v", err)
		return false, 0
	}
	if len(res) > 0 {
		return true, res[0].Id
	}
	return false, 0
}

func (t *PlaybookService) SubmitPlayBook(ctx context.Context, reqx *proto.SubmitPlayBookReq) (*proto.SubmitPlayBookResponse, error) {
	if t.Limit != nil && !t.Limit.Request() {
		return nil, errors.New(fmt.Sprintf("%v qps limitation is: %d", t.Name, t.Limit.GetQuota()))
	}

	dao := orm.GetEngine()
	var reqjson map[string]interface{}
	err := json.Unmarshal(reqx.Body, &reqjson)
	if err != nil {
		return nil, err
	}

	jsonb, _ := json.Marshal(reqjson)

	var req core.RawPlaybook
	err = json.Unmarshal(jsonb, &req)
	if err != nil {
		return nil, err
	}

	var pbs = make([]models.Playbook, 0)
	err = dao.Slave().Where("id = ?", req.PlaybookId).Find(&pbs)
	if err != nil {
		logx.L().Errorf("find err:%v", err)
		return nil, err
	}
	if len(pbs) != 1 {
		return nil, fmt.Errorf("pbid not exist:%d", req.PlaybookId)
	}

	// 鉴权
	md, _ := metadata.FromIncomingContext(ctx)
	uif, ok := md["auth"]
	if !ok || len(uif) != 1 {
		return nil, fmt.Errorf("FromIncomingContext has no auth")
	}
	var user auth.UserInfo
	err = json.Unmarshal(toolutil.String2Byte(uif[0]), &user)
	if err != nil {
		return nil, err
	}

	err = user.HasAppPermit(pbs[0].AppId, true)
	if err != nil {
		return nil, err
	}

	// 检验剧本合法性
	pb, err := core.NewPlayBook(bytes.NewBuffer(jsonb))
	if err != nil {
		return nil, err
	}
	snapb, _ := json.Marshal(pb)
	sum := toolutil.CheckSum(snapb)
	exist, sid := CheckSumExist(sum)

	session := dao.Master().NewSession()
	defer session.Close()
	_ = session.Begin()

	if !exist {
		// 创建新的快照
		newsnap := models.Snapshot{
			PlaybookId: pb.GetId(),
			Snapshot:   toolutil.Bytes2string(snapb),
			Rawbody:    toolutil.Bytes2string(jsonb),
			Checksum:   sum,
			AppId:      pbs[0].AppId,
			Snapname:   fmt.Sprintf("snapshot-%s", sum),
			UpdateTime: time.Now(),
			User:       user.Account,
		}
		//dao.Master().Insert(&newsnap)
		_, err = session.Insert(&newsnap)
		if err != nil {
			_ = session.Rollback()
			return nil, err
		}

		// 创建完快照后，更新剧本快照id
		updatepb := models.Playbook{
			SnapshotId: newsnap.Id,
			UpdateTime: time.Now(),
		}
		//dao.Master().Where("id = ?", req.PlaybookId).Update(&updatepb)
		_, err = session.Where("id = ?", req.PlaybookId).Update(&updatepb)
		if err != nil {
			_ = session.Rollback()
			return nil, err
		}

		// 通知集群，后续做内存调整，详情查看dcc模块
		err = dcc.GetDcc().DccPutPb(req.PlaybookId, newsnap.Id)
		if err != nil {
			_ = session.Rollback()
			return nil, err
		}

	} else {

		// 直接更新raw body
		update := models.Snapshot{
			Rawbody:    toolutil.Bytes2string(jsonb),
			UpdateTime: time.Now(),
		}
		//dao.Master().Where("checksum = ? ", sum).Update(&update)
		_, err = session.Where("checksum = ? ", sum).Update(&update)
		if err != nil {
			_ = session.Rollback()
			return nil, err
		}

		// 更新剧本的快照id
		updatepb := models.Playbook{
			SnapshotId: sid,
			UpdateTime: time.Now(),
		}
		//dao.Master().Where("id = ?", req.PlaybookId).Update(&updatepb)
		_, err = session.Where("id = ?", req.PlaybookId).Update(&updatepb)
		if err != nil {
			_ = session.Rollback()
			return nil, err
		}

		// 相当于做了一次版本切换
		err = dcc.GetDcc().DccPutPb(req.PlaybookId, sid)
		if err != nil {
			logx.L().Errorf("err when put etcd", err)
			_ = session.Rollback()
			return nil, err
		}
	}

	_ = session.Commit()

	return &proto.SubmitPlayBookResponse{
		PlaybookId: int32(pb.GetId()),
	}, nil
}

func (t *PlaybookService) GetPlaybook(ctx context.Context, req *proto.GetPlaybookReq) (*proto.GetPlaybookResponse, error) {
	pbcache, ok := core.G().GetPlaybook(cast.ToInt(req.Pbid))
	if !ok {
		return nil, fmt.Errorf("pb not exist")
	}
	resp := map[string]interface{}{
		"id":    pbcache.GetId(),
		"nodes": pbcache.GetNodes(),
	}
	b, _ := json.Marshal(resp)
	return &proto.GetPlaybookResponse{
		Data: b,
	}, nil
}

func (t *PlaybookService) SwitchSnapshot(ctx context.Context, req *proto.SwitchVersionReq) (*proto.SwitchVersionResponse, error) {
	if t.Limit != nil && !t.Limit.Request() {
		return nil, errors.New(fmt.Sprintf("%v qps limitation is: %d", t.Name, t.Limit.GetQuota()))
	}
	dao := orm.GetEngine()
	// 切换完快照后，更新剧本对应的快照id
	updatepb := models.Playbook{
		SnapshotId: int(req.SnapShotId),
		UpdateTime: time.Now(),
	}
	_, err := dao.Master().Where("id = ?", req.PlayBookId).Update(&updatepb)
	if err != nil {
		return nil, err
	}

	// 通知集群，后续做内存调整，详情查看dcc模块
	err = dcc.GetDcc().DccPutPb(int(req.PlayBookId), int(req.SnapShotId))
	if err != nil {
		return nil, err
	}
	return &proto.SwitchVersionResponse{
		PlayBookId: req.PlayBookId,
	}, nil
}
