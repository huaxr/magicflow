package logx

import (
	"io"
	"os"
	"time"

	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/huaxr/magicflow/pkg/toolutil"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger Logger

type Logger struct {
	levelLogger *zap.SugaredLogger
}

func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.FullPath() + "]")
}

func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func init() {
	var encoder zapcore.Encoder
	switch confutil.GetLog().GetEncoder() {
	case confutil.Json:
		// {"level":"INFO","time":"2022-04-06 14:05:31","file":"consensus/etcdcli.go:78","msg":"Elect success","host":"10.74.152.206"}
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "file",
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeTime:   cEncodeTime,
		})

	case confutil.Console:
		// 2022-04-06 11:43:12	[INFO]	[publisher/nsq.go:131]	start/recover 10.90.72.172:4150 nsq publisher
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "file",
			EncodeCaller: zapcore.FullCallerEncoder,
			EncodeLevel:  zapcore.CapitalColorLevelEncoder,
			EncodeTime:   cEncodeTime,
		})
	}

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= confutil.GetLog().GetLogLevel()
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	infoWriter := getWriter(confutil.GetLog().Infofile)
	errorWriter := getWriter(confutil.GetLog().Errorfile)

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	}

	if confutil.GetLog().GetConsole() {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), infoLevel))
	}
	core := zapcore.NewTee(cores...)
	log := zap.New(core, zap.AddCaller())
	logger = Logger{
		log.Sugar(),
	}
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+"%Y%m%d"+".log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

// tags； zap.String("metric", "demoApp"), zap.String("tag1", "demoApp")
func L(fields ...interface{}) *zap.SugaredLogger {
	if confutil.GetLog().GetDisabletags() {
		return logger.levelLogger
	}

	if len(fields) > 0 {
		return logger.levelLogger.With(zap.String("host", toolutil.GetIp())).With(fields...)
	}

	return logger.levelLogger.With(zap.String("host", toolutil.GetIp())).With()

}
