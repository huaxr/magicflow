// Author: XinRui Hua
// Time:   2022/4/12 下午2:18
// Git:    huaxr

package orm

import (
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/pkg/confutil"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var (
	eg *xorm.EngineGroup

	slowDuration = 1000 * time.Millisecond
	showSql      = false
	maxConn      = 100 //Maximum number of connections
	maxIdle      = 30  //Maximum number of idle connections
	logLevel     = log.LOG_INFO
)

func initvar() {
	if confutil.GetConf().Db.Mysql.LogLevel >= 0 {
		logLevel = log.LogLevel(confutil.GetConf().Db.Mysql.LogLevel)
	}

	if confutil.GetConf().Db.Mysql.ShowSql {
		showSql = true
	}

	if confutil.GetConf().Db.Mysql.MaxConn > 0 {
		maxConn = confutil.GetConf().Db.Mysql.MaxConn
	}

	if confutil.GetConf().Db.Mysql.MaxIdle > 0 {
		maxIdle = confutil.GetConf().Db.Mysql.MaxIdle
	}

	if confutil.GetConf().Db.Mysql.SlowDuration > 0 {
		slowDuration = time.Duration(confutil.GetConf().Db.Mysql.SlowDuration) * time.Millisecond
	}
}

func LaunchDbEngine() {
	initvar()
	var err error
	master, err := xorm.NewEngine("mysql", confutil.GetConf().Db.Mysql.Master)
	if err != nil {
		logx.L().Errorf("NewEngine err: %v", err)
		return
	}

	slaves := make([]*xorm.Engine, 0)
	for _, slave := range confutil.GetConf().Db.Mysql.Slaves {
		s, err := xorm.NewEngine("mysql", slave)
		if err != nil {
			logx.L().Errorf("NewEngine err: %v with slave: %v", err, slave)
			return
		}
		slaves = append(slaves, s)
	}

	eg, err = xorm.NewEngineGroup(master, slaves)
	if err != nil {
		logx.L().Panicf("err %v", err)
		return
	}

	if err = eg.Ping(); err != nil {
		logx.L().Panicf("ping err %v", err)
	}

	eg.SetMaxIdleConns(maxIdle)
	eg.SetMaxOpenConns(maxConn)
	eg.SetLogger(dbLogger)
}

func GetEngine() *xorm.EngineGroup {
	return eg
}
