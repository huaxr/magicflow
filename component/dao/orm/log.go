// Author: XinRui Hua
// Time:   2022/4/13 上午9:58
// Git:    huaxr

package orm

import (
	"fmt"
	"time"

	"github.com/huaxr/magicflow/component/logx"
	"xorm.io/xorm/log"
)

var dbLogger = &dbLog{}

type dbLog struct{}

func (d *dbLog) Debug(v ...interface{}) {
	logx.L().Debug(fmt.Sprint(v...))
}

func (d *dbLog) Debugf(format string, v ...interface{}) {
	logx.L().Debugf("sql-debug:%v", fmt.Sprint(v...))
}

func (d *dbLog) Info(v ...interface{}) {
	logx.L().Info(fmt.Sprint(v...))
}

func (d *dbLog) Infof(format string, v ...interface{}) {
	if len(v) > 0 {
		duration, ok := v[len(v)-1].(time.Duration)
		if ok && duration < slowDuration {
			return
		}
	}

	logx.L().Debugf("multi insert slow sql")
}

func (d *dbLog) Warn(v ...interface{}) {
	logx.L().Warn("sql-warn", fmt.Sprint(v...))
}

func (d *dbLog) Warnf(format string, v ...interface{}) {
	logx.L().Warnf(format, fmt.Sprint(v...))

}

func (d *dbLog) Error(v ...interface{}) {
	logx.L().Error("sql-error", fmt.Sprint(v...))
}

func (d *dbLog) Errorf(format string, v ...interface{}) {
	logx.L().Errorf(format, fmt.Sprint(v...))
}

func (d *dbLog) Level() log.LogLevel {
	return log.LOG_INFO
}

func (d *dbLog) SetLevel(l log.LogLevel) {
	return
}

func (d *dbLog) ShowSQL(show ...bool) {
	return
}

func (d *dbLog) IsShowSQL() bool {
	return showSql
}
