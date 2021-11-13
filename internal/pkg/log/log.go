package log

import (
	"bytes"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/grpclog"
	"io"
	"os"
)

func init() {
	New(os.Stderr, false)
}

type logFunc2 func(template string, args ...interface{})
type logFunc func(args ...interface{})

func (the logFunc) Print(args ...interface{}) {
	the(args)
}

var (
	Info  logFunc
	Warn  logFunc
	Error logFunc
	Debug logFunc
	Panic logFunc
	Fatal logFunc

	Infof  logFunc2
	Warnf  logFunc2
	Errorf logFunc2
	Debugf logFunc2
	Panicf logFunc2
)

var Z *zap.Logger
var S *zap.SugaredLogger

func New(writer io.Writer, isProduction bool) {
	Z = newLogger(writer, isProduction)
	S = Z.Sugar()

	Info = S.Info
	Warn = S.Warn
	Error = S.Error
	Debug = S.Debug
	Panic = S.Panic

	Infof = S.Infof
	Warnf = S.Warnf
	Errorf = S.Errorf
	Debugf = S.Debugf
	Panicf = S.Panicf
	Fatal = S.Fatal
}

func newDevelopEncoder() zapcore.Encoder {
	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.EncodeCaller = zapcore.FullCallerEncoder
	return zapcore.NewConsoleEncoder(encCfg)
}

func newLogger(writer io.Writer, isProduction bool) *zap.Logger {
	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.DebugLevel
	})
	enc := newDevelopEncoder()
	if isProduction {
		lowPriority = func(lv zapcore.Level) bool { return lv >= zapcore.InfoLevel }
		enc = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	// addSync 将 io.Writer 装饰为 WriteSyncer
	// 故只需要一个实现 io.Writer 接口的对象即可
	redisCore := zapcore.NewCore(enc, zapcore.AddSync(writer), lowPriority)
	core := zapcore.NewTee(redisCore)
	// 集成多个 ddd
	// 使用 JSON 格式日志
	stdCore := zapcore.NewCore(enc, zapcore.Lock(os.Stdout), lowPriority)
	if isProduction {
		lowPriority = func(lv zapcore.Level) bool { return lv >= zapcore.WarnLevel }
		stdCore = zapcore.NewCore(enc, zapcore.Lock(os.Stderr), lowPriority)
	}
	core = zapcore.NewTee(redisCore, stdCore)

	// logger 输出到 console 且标识调用代码行
	return zap.New(core).WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func LoggerGRpc() grpclog.LoggerV2 {
	buf := bytes.NewBuffer(nil)
	lg := zap.New(Z.Core(), zap.AddCaller(), zap.AddCallerSkip(2), zap.ErrorOutput(zapcore.AddSync(buf)))
	return &zapGRPCLogger{lg: lg, sugar: lg.Sugar()}
}

func ErrorLogger() mysql.Logger {
	return Error
}
