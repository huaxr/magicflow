package playbookservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/logx"

	"time"

	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/dcc"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/gin-gonic/gin"
)

func CreateEmptyPlayBook(c *gin.Context, req *CreatePlayBookReq) (CreatePlayBookResp, error) {
	// 鉴权
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err := user.HasAppPermit(req.AppId, true)
	if err != nil {
		return CreatePlayBookResp{}, err
	}

	dao := orm.GetEngine()
	// 检查app是否存在
	var app = make([]models.App, 0)
	dao.Slave().Where("id = ?", req.AppId).Find(&app)
	if len(app) == 0 {
		return CreatePlayBookResp{}, errors.New(fmt.Sprintf("app %v not exist", req.AppId))
	}

	// Check the uniqueness of the name, name should be unique in the same app
	var res = make([]models.Playbook, 0)
	err = dao.Slave().Where("app_id = ? and name = ?", req.AppId, req.Name).Find(&res)
	if err != nil {
		return CreatePlayBookResp{}, err
	}
	if len(res) != 0 {
		return CreatePlayBookResp{}, errors.New("name already exists")
	}

	// Create new playbook
	create := models.Playbook{
		AppId:       req.AppId,
		Name:        req.Name,
		Enable:      1,
		Description: req.Description,
		UpdateTime:  time.Now(),
		User:        user.Account,
		Token:       toolutil.GetRandomString(8),
	}

	_, err = dao.Master().Insert(&create)
	if err != nil {
		return CreatePlayBookResp{}, err
	}
	id := create.Id

	// generate related task automatically
	insert := models.Tasks{
		Configuration: "{}",
		Name:          fmt.Sprintf("剧本任务:%v", id),
		Description:   "自动生成的本地剧本任务",
		AppId:         req.AppId,
		Xrn:           core.NewLocalTaskRN(id),
		Type:          "local",
		UpdateTime:    time.Now(),
		InputExample:  "",
		OutputExample: "",
		User:          user.Account,
	}
	dao.Master().Insert(&insert)

	return CreatePlayBookResp{id}, nil
}

// CheckSum function checks whether body has been loaded
// true means body has been loaded, we should refuse
func CheckSumExist(sum string) (bool, int) {
	dao := orm.GetEngine()
	var res = make([]models.Snapshot, 0)
	dao.Slave().Where("checksum = ?", sum).Find(&res)
	if len(res) > 0 {
		return true, res[0].Id
	}
	return false, 0
}

func SubmitPlayBook(c *gin.Context, reqjson map[string]interface{}) (CreatePlayBookResp, error) {
	dao := orm.GetEngine()
	jsonb, _ := json.Marshal(reqjson)

	var req core.RawPlaybook
	err := json.Unmarshal(jsonb, &req)
	if err != nil {
		return CreatePlayBookResp{}, err
	}

	var pbs = make([]models.Playbook, 0)
	dao.Slave().Where("id = ?", req.PlaybookId).Find(&pbs)
	if len(pbs) != 1 {
		return CreatePlayBookResp{}, errors.New(fmt.Sprintf("pbid not exist:%d", req.PlaybookId))
	}

	// 鉴权
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err = user.HasAppPermit(pbs[0].AppId, true)
	if err != nil {
		return CreatePlayBookResp{}, err
	}

	// 检验剧本合法性
	pb, err := core.NewPlayBook(bytes.NewBuffer(jsonb))
	if err != nil {
		return CreatePlayBookResp{}, err
	}
	snapb, _ := json.Marshal(pb)
	sum := toolutil.CheckSum(snapb)
	exist, sid := CheckSumExist(sum)

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
		dao.Master().Insert(&newsnap)

		// 创建完快照后，更新剧本快照id
		updatepb := models.Playbook{
			SnapshotId: newsnap.Id,
			UpdateTime: time.Now(),
		}
		dao.Master().Where("id = ?", req.PlaybookId).Update(&updatepb)

		// 通知集群，后续做内存调整，详情查看dcc模块
		dcc.GetDcc().DccPutPb(req.PlaybookId, newsnap.Id)
	} else {

		// 直接更新raw body
		update := models.Snapshot{
			Rawbody:    toolutil.Bytes2string(jsonb),
			UpdateTime: time.Now(),
		}
		dao.Master().Where("checksum = ? ", sum).Update(&update)

		// 更新剧本的快照id
		updatepb := models.Playbook{
			SnapshotId: sid,
			UpdateTime: time.Now(),
		}
		dao.Master().Where("id = ?", req.PlaybookId).Update(&updatepb)

		// 相当于做了一次版本切换
		dcc.GetDcc().DccPutPb(req.PlaybookId, sid)
	}
	return CreatePlayBookResp{pb.GetId()}, nil
}

func GetRealPlayBook(ctx context.Context, req *GetPlayBookReq) (GetPlayBookResp, error) {
	dao := orm.GetEngine()

	var res = make([]models.Playbook, 0)
	err := dao.Slave().Where("id = ?", req.PlayBookId).Find(&res)
	if err != nil {
		logx.L().Errorf("GetRealPlayBook find by name is err: %+v", err)
		return GetPlayBookResp{}, err
	}
	if len(res) != 1 {
		return GetPlayBookResp{}, errors.New("no playbook found")
	}
	return GetPlayBookResp{res[0]}, nil
}

func SwitchSnapshot(c *gin.Context, req *SwitchVersionReq) (SwitchVersionResp, error) {
	dao := orm.GetEngine()

	// 切换完快照后，更新剧本对应的快照id
	updatepb := models.Playbook{
		SnapshotId: req.SnapShotId,
		UpdateTime: time.Now(),
	}
	dao.Master().Where("id = ?", req.PlayBookId).Update(&updatepb)

	// 通知集群，后续做内存调整，详情查看dcc模块
	err := dcc.GetDcc().DccPutPb(req.PlayBookId, req.SnapShotId)
	if err != nil {
		return SwitchVersionResp{}, err
	}
	return SwitchVersionResp{req.PlayBookId}, nil
}
