package playbookcontroller

import (
	"encoding/json"
	"fmt"

	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/logx"

	"net/http"

	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/api/controller"
	"github.com/huaxr/magicflow/component/api/service/playbookservice"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type PlayBookController struct {
	controller.BaseApiController
}

func (ctl *PlayBookController) Router(router *gin.Engine) {
	entry := router.Group("/playbook")
	//// 获取剧本后端数据结构详情
	//entry.GET("/get_real_playbook", ctl.GetRealPlayBook)
	// 获取剧本前端数据结构详情
	entry.GET("/get_raw_playbook", ctl.GetRawPlaybook)

	// 根据 trace_id 获取执行列表
	entry.GET("/get_trace_detail", ctl.GetTrace)
	// 根据 snapshot_id 获取执行列表（考虑到执行结果需要和快照id的绑定渲染不同的DGA）
	entry.GET("/get_trace_list", ctl.GetTraceBySnapshotId)

	// 获取快照列表
	entry.GET("/get_snapshots", ctl.GetSnapshots) // snapshot list
	// 获取快照详情
	entry.GET("/get_snapshot_detail", ctl.GetSnapshotDetail) // snapshot
	//// 切换快照
	//entry.POST("/switch_snapshot", ctl.SwitchSnapshot)

	// 创建空剧本
	entry.POST("/create", ctl.CreateEmptyPlayBook)
	//// 提交剧本并保存 raw 和 real struct，并通知dcc
	//entry.POST("/submit", ctl.SubmitPlaybook)

	// ugly demo api
	//entry.GET("/get_execution", ctl.GetExecution)
	//entry.GET("/get_es_execution", ctl.GetEsExecution)
	//entry.GET("/get_context", ctl.GetContext)
	//entry.GET("/get_es_context", ctl.GetEsContext)
}

// get statemachine information from pbcache and db
func (ctl *PlayBookController) GetRealPlayBook(c *gin.Context) {
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
	pbcache, ok := core.G().GetPlaybook(cast.ToInt(smId))
	if !ok {
		ctl.Error(c, 200, "pb not exist", nil)
		return
	}
	resp := map[string]interface{}{
		"id":          pbcache.GetId(),
		"nodes":       pbcache.GetNodes(),
		"name":        pb.Name,
		"description": pb.Description,
		"snapshot_id": pb.SnapshotId,
	}
	ctl.Success(c, resp)
}

func (ctl *PlayBookController) GetRawPlaybook(c *gin.Context) {
	pbId, exist := c.GetQuery("id")
	if !exist || cast.ToInt(pbId) <= 0 {
		logx.L().Errorf("id: %v not right!", pbId)
		ctl.Error(c, 200, "err id", nil)
		return
	}

	dao := orm.GetEngine()

	var pbs = make([]models.Playbook, 0)
	// 获取当前启用的剧本快照
	err := dao.Slave().Where("id = ?", pbId).Find(&pbs)
	if err != nil {
		ctl.Error(c, 200, "playbook query err", err)
		return
	}

	if len(pbs) == 0 {
		ctl.Error(c, 200, "playbook not exist", err)
		return
	}

	var ss = make([]models.Snapshot, 0)
	// 获取当前启用的剧本快照
	err = dao.Slave().Where("id = ?", pbs[0].SnapshotId).Find(&ss)
	if err != nil {
		logx.L().Errorf("GetRawPlaybook err %v", err)
		ctl.Error(c, 200, "GetRawPlaybook is err", err)
		return
	}

	// 如果是空的话，代表还没有剧本, 直接反回
	if len(ss) == 0 {
		ctl.Success(c, gin.H{})
		return
	}

	var resp map[string]interface{}
	json.Unmarshal([]byte(ss[0].Rawbody), &resp)
	ctl.Success(c, resp)
}

type Connections struct {
	From string
	To   string
	Desc string
	Exec bool
}

func getIndex(l []string, v string) int {
	for index, value := range l {
		if value == v {
			return index
		}
	}
	return -1
}
func draw(path []string, from, to string) bool {
	i1 := getIndex(path, from)
	i2 := getIndex(path, to)
	return i2-i1 == 1
}

func (ctl *PlayBookController) GetTrace(c *gin.Context) {
	exId, ok := c.GetQuery("trace_id")
	if !ok {
		ctl.Error(c, 200, "trace_id is empty", nil)
		return
	}
	ssId, ok := c.GetQuery("snapshot_id")
	if !ok {
		ctl.Error(c, 200, "snapshot_id is empty", nil)
		return
	}

	dao := orm.GetEngine()
	var exe = make([]models.Execution, 0)
	table := fmt.Sprintf("execution_%d", cast.ToInt(ssId)%10)
	dao.Slave().Table(table).Where("trace_id = ? and domain = '' and node_code != ''", exId).OrderBy("sequence").Find(&exe)
	ctl.Success(c, exe)
}

func (ctl *PlayBookController) GetTraceBySnapshotId(c *gin.Context) {
	ssId, _ := c.GetQuery("snapshot_id")

	status := c.DefaultQuery("status", "success")
	dao := orm.GetEngine()
	offset, limit := toolutil.GetOffsetLimit(c)

	//offset, limit := toolutil.GetOffsetLimit(c)
	var exe = make([]models.Execution, 0)
	table := fmt.Sprintf("execution_%d", cast.ToInt(ssId)%10)

	switch status {
	case "fail":
		dao.Slave().Table(table).Where("snapshot_id = ? and status = ?", ssId, "TaskFail").Limit(limit, offset).Cols("trace_id").OrderBy("id desc").Distinct("trace_id").Find(&exe)
	default:
		dao.Slave().Table(table).Where("snapshot_id = ? and node_code = ?", ssId, "").Limit(limit, offset).Cols("trace_id").OrderBy("id desc").Find(&exe)
	}
	var traceList = make([]string, 0)

	for _, i := range exe {
		traceList = append(traceList, i.TraceId)
	}

	ctl.Success(c, traceList)
}

func (ctl *PlayBookController) GetSnapshots(c *gin.Context) {
	pbId, exist := c.GetQuery("playbook_id")
	if !exist || cast.ToInt(pbId) <= 0 {
		logx.L().Errorf("playbook_id: %v not right!", pbId)
		ctl.Error(c, 200, "err playbook_id", nil)
		return
	}
	dao := orm.GetEngine()
	var pbs = make([]models.Playbook, 0)
	// 获取当前启用的剧本快照
	err := dao.Slave().Where("id = ?", pbId).Find(&pbs)
	if err != nil {
		ctl.Error(c, 200, "playbook query err", err)
		return
	}

	if len(pbs) == 0 {
		ctl.Error(c, 200, "playbook not exist", nil)
		return
	}

	// 鉴权
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err = user.HasAppPermit(cast.ToInt(pbs[0].AppId), true)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	var snaps []models.Snapshot
	dao.Slave().Where("playbook_id = ?", pbId).Cols("id, playbook_id, app_id, checksum, snapname, update_time").OrderBy("id desc").Limit(10).Find(&snaps)

	ctl.Success(c, gin.H{
		"snapshots_list":      snaps,
		"current_snapshot_id": pbs[0].SnapshotId,
	})
	return
}

func (ctl *PlayBookController) GetSnapshotDetail(c *gin.Context) {

	pbid, exist := c.GetQuery("playbook_id")
	if !exist || cast.ToInt(pbid) <= 0 {
		logx.L().Errorf("playbook_id: %v not right!", pbid)
		ctl.Error(c, 200, "err playbook_id", nil)
		return
	}

	ssid, exist := c.GetQuery("snapshot_id")
	if !exist || cast.ToInt(ssid) <= 0 {
		ctl.Success(c, gin.H{})
		return
	}

	dao := orm.GetEngine()

	var snaps []models.Snapshot
	dao.Slave().Where("id = ? and playbook_id = ?", ssid, pbid).Find(&snaps)

	ctl.Success(c, snaps[0])
	return
}

func (ctl *PlayBookController) SwitchSnapshot(c *gin.Context) {
	req := &playbookservice.SwitchVersionReq{}
	if err := c.ShouldBind(req); err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	dao := orm.GetEngine()

	var res = make([]models.Snapshot, 0)
	dao.Slave().Where("id = ? and playbook_id = ?", req.SnapShotId, req.PlayBookId).Find(&res)

	if len(res) != 1 {
		ctl.Error(c, 200, fmt.Sprintf("snapshot %v not exist", req.SnapShotId), nil)
		return
	}
	// 鉴权
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err := user.HasAppPermit(cast.ToInt(res[0].AppId), true)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	rep, err := playbookservice.SwitchSnapshot(c, req)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	ctl.Success(c, rep)
}

func (ctl *PlayBookController) GetExecution(c *gin.Context) {
	var nodes = make([]string, 0)
	var connections = make([]Connections, 0)
	var executed = make([]string, 0)
	var success = make([]string, 0)
	var fail = make([]string, 0)

	exId, _ := c.GetQuery("trace_id")
	dao := orm.GetEngine()

	var exe []models.Execution
	dao.Slave().Where("trace_id = ? and node_code != ? and domain = ?", exId, core.StartNodeCode, "").OrderBy("sequence").Find(&exe)

	if len(exe) == 0 {
		ctl.Error(c, 200, "query not exist", nil)
		return
	}

	var path = make([]string, 0)

	for _, i := range exe {
		path = append(path, i.NodeCode)
		switch i.Status {
		case "TaskExecuted":
			executed = append(executed, i.NodeCode)
		case "TaskSuccess":
			success = append(success, i.NodeCode)
		default:
			fail = append(fail, i.NodeCode)
		}
	}

	pbs := core.G()
	pb, ok := pbs.GetPlaybook(exe[0].PlaybookId)
	if !ok {
		c.JSON(200, map[string]interface{}{"err": "pb is nil"})
		return
	}
	for k, v := range pb.GetNodes() {
		nodes = append(nodes, k)
		for _, e := range v.GetBranchChoice() {
			// In order to determine whether this road is strong implementation
			var exec1, exec2 bool
			var ex []models.Execution
			dao.Slave().Where("trace_id = ? and node_code = ? and domain = ?", exId, e.NextNode, "").Find(&ex)
			if len(ex) == 1 {
				exec1 = true
			}

			ex = ex[:0]
			dao.Slave().Where("trace_id = ? and node_code = ? and domain = ?", exId, v.GetNodeCode(), "").Find(&ex)
			if len(ex) == 1 {
				exec2 = true
			}

			c := Connections{
				From: k,
				To:   e.NextNode,
				Exec: exec1 && exec2,
			}
			connections = append(connections, c)
		}

		for _, e := range v.GetBranchParallel() {
			c := Connections{
				From: k,
				To:   e.NextNode,
				Exec: true,
			}
			connections = append(connections, c)
		}
	}

	var file string = "index.html"
	c.HTML(http.StatusOK, file, gin.H{
		"title":       "编排项目",
		"nodes":       nodes,
		"connections": connections,
		"executed":    executed,
		"success":     success,
		"fail":        fail,
	})
}

func (ctl *PlayBookController) GetContext(c *gin.Context) {
	exId, _ := c.GetQuery("trace_id")
	nodeName, _ := c.GetQuery("node_name")
	dao := orm.GetEngine()

	var exe []models.Execution
	_ = dao.Slave().Where("trace_id = ? and node_code = ?", exId, nodeName).Find(&exe)

	if len(exe) != 1 {
		ctl.Error(c, 200, "not unique node", nil)
		return
	}

	var response map[string]interface{}
	b, _ := json.Marshal(exe[0])
	_ = json.Unmarshal(b, &response)

	ok := core.G().HasPlaybook(exe[0].PlaybookId)
	if !ok {
		ctl.Error(c, 200, "pb is nil", nil)
		return
	}
	ctl.Success(c, response)
}

func (ctl *PlayBookController) CreateEmptyPlayBook(c *gin.Context) {
	var params playbookservice.CreatePlayBookReq
	err := c.ShouldBind(&params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ret, err := playbookservice.CreateEmptyPlayBook(c, &params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ctl.Success(c, ret)
	return
}

func (ctl *PlayBookController) SubmitPlaybook(c *gin.Context) {
	var reqjson map[string]interface{}
	err := c.ShouldBind(&reqjson)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ret, err := playbookservice.SubmitPlayBook(c, reqjson)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ctl.Success(c, ret)
	return
}
