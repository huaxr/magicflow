// Author: huaxr
// Time:   2021/6/28 下午3:32
// Git:    huaxr

package core

import (
	ctx2 "context"
	"encoding/json"
	"fmt"

	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/ticker"
	"xorm.io/xorm"

	"runtime"
	"time"

	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/spf13/cast"
)

const (
	tableSlot = 10
	multi     = 500

	flushTick = 1000 * time.Millisecond
	delay     = 2022 * time.Millisecond
)

func wait() {
	time.Sleep(flushTick + delay)
	runtime.Gosched()
}

type recorder struct {
	dao           *xorm.EngineGroup
	ctx           ctx2.Context
	executionChan chan *models.Execution
	inserting     [tableSlot]*[]*models.Execution
}

func initDaoJob(ctx ctx2.Context) *recorder {
	dj := new(recorder)
	dao := orm.GetEngine()
	dj.dao = dao
	dj.ctx = ctx
	dj.executionChan = make(chan *models.Execution, 1e7)

	for i := 0; i < tableSlot; i++ {
		exe := make([]*models.Execution, 0)
		dj.inserting[i] = &exe
	}

	return dj
}

func getTableSlot(ssId int) int {
	return ssId % tableSlot
}

// recordImmediately
func (d *recorder) recordImmediately() {
	for index, pendInsert := range d.inserting {
		if len(*pendInsert) == 0 {
			continue
		}
		var size = len(*pendInsert)
		if size >= multi {
			size = multi
		}

		logx.L().Debugf("flush %v record into %v", size, fmt.Sprintf("execution_%d", index))
		_, err := d.dao.Master().Table(fmt.Sprintf("execution_%d", index)).Insert((*pendInsert)[:size])
		if err != nil {
			// invalid connection maybe
			logx.L().Errorf("recordDbImmediately err:%v", err)
			continue
		}

		*pendInsert = (*pendInsert)[size:]
	}
}

type Extra struct {
	Input     interface{} `json:"input,omitempty"`
	Output    interface{} `json:"output,omitempty"`
	Exception interface{} `json:"exception,omitempty"`
	Detail    *message    `json:"detail,omitempty"`
}

func (d *recorder) Name() string           { return "db_record" }
func (d *recorder) Duration() *time.Ticker { return time.NewTicker(flushTick) }
func (d *recorder) Heartbeat() {
	ticker.GetManager().Register(ticker.NewJob(d.ctx, d.Name(), d.Duration(), d.recordImmediately))
}

// startRecord one goroutine
// table strategy: range by slot or delivery (24) by app_name
func (d *recorder) startRecord() {
	ticker.RegisterTick(d)
	for {
		select {
		case execute := <-d.executionChan:
			s := getTableSlot(execute.SnapshotId)
			*d.inserting[s] = append(*d.inserting[s], execute)

		case <-d.ctx.Done():
			// record lag behind messages interim the os.exit waiting period which
			// due to eliminating slot lapse. surplus recorder may cause unexpected
			// lost.
			d.recordImmediately()
			logx.L().Warnf("record thread shutdown")
			return
		}
	}
}

func loadoutputCtxFromBackend(ssId int, traceId int64, nodeCode string) interface{} {
	logx.L().Infof("query data from db by trace_id:%v", traceId)
	dao := orm.GetEngine()
	var res = make([]models.Execution, 0)
	table := fmt.Sprintf("execution_%d", cast.ToInt(ssId)%10)
	dao.Slave().Table(table).Where("trace_id = ? and node_code = ?", cast.ToString(traceId), nodeCode).Find(&res)

	if len(res) == 0 {
		return nil
	}

	var extra Extra
	err := json.Unmarshal([]byte(res[0].Extra), &extra)
	if err != nil {
		logx.L().Errorf("loadoutputCtxFromBackend %v", err)
		return nil
	}
	return extra.Output
}

func loadFirst() {
	dao := orm.GetEngine()

	var res = make([]models.Playbook, 0)
	err := dao.Slave().Where("enable = ? and snapshot_id > 0", 1).Find(&res)
	if err != nil {
		logx.L().Errorf("loadFirst err %v", err)
		return
	}

	logx.L().Infof("loadFirst total get:%d pb", len(res))
	for _, pb := range res {
		err = ReloadSnapshot(pb.SnapshotId, pb.Id)
		if err != nil {
			logx.L().Errorf("loadFirst err:%v, bypass:%v", err, pb.Id)
			continue
		}
	}

	var apps = make([]models.App, 0)
	err = dao.Slave().Where("checked = ? and brokers != ?", 1, "").Find(&apps)
	if err != nil {
		logx.L().Errorf("loadapp err %v", err)
		return
	}

	for _, i := range apps {
		SetNamespace(&i)
	}
}

func (m *message) UpdateStatus(stat TaskStatus) {
	logx.L().Debugf("update execution ")
	dao := orm.GetEngine()
	ex := models.Execution{
		Status: string(stat),
	}
	table := fmt.Sprintf("execution_%d", cast.ToInt(m.Task.SnapshotId)%10)
	_, err := dao.Master().Table(table).Where("trace_id = ? and node_code = ? and domain = ?", m.Meta.Trace, m.Task.NodeCode, m.Meta.Domain).Update(&ex)
	if err != nil {
		logx.L().Errorf("update err %v", err)
	}
}

func findPb(pbid int) (*models.Playbook, error) {
	// reload to check the "real" snapshot.
	dao := orm.GetEngine()
	var pb = new(models.Playbook)
	_, err := dao.Slave().ID(pbid).Get(pb)
	if err != nil {
		return nil, err
	}

	return pb, nil
}

func findApp(appid int) (*models.App, error) {
	dao := orm.GetEngine()
	var app = new(models.App)
	_, err := dao.Slave().ID(appid).Get(app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func findTask(taskid int) (*models.Tasks, error) {
	dao := orm.GetEngine()
	var task = new(models.Tasks)
	_, err := dao.Slave().ID(taskid).Get(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func findSnap(snapId int) (*models.Snapshot, error) {
	dao := orm.GetEngine()
	var ss = new(models.Snapshot)
	_, err := dao.ID(snapId).Get(ss)
	if err != nil {
		return nil, err
	}

	if ss.Id == 0 {
		return nil, fmt.Errorf("snapshot not exist")
	}
	return ss, nil
}
