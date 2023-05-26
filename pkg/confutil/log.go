// Author: XinRui Hua
// Time:   2022/4/5 上午11:05
// Git:    huaxr

package confutil

import "go.uber.org/zap/zapcore"

type EncodeType string

const (
	Console EncodeType = "console"
	Json    EncodeType = "json"
)

type Log struct {
	Level       string     `yaml:"level"`
	Encoder     EncodeType `yaml:"encoder"`
	Console     bool       `yaml:"console"`
	Disabletags bool       `yaml:"disabletags"`
	Infofile    string     `yaml:"infofile"`
	Errorfile   string     `yaml:"errorfile"`
}

var log *Log

func GetLog() *Log {
	if log == nil {
		initConf()
	}
	return log
}

func (l *Log) GetLogLevel() zapcore.Level {
	switch l.Level {
	default:
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}
}

func (l *Log) GetEncoder() EncodeType {
	return l.Encoder
}

func (l *Log) GetConsole() bool {
	return l.Console
}

func (l *Log) GetDisabletags() bool {
	return l.Disabletags
}
