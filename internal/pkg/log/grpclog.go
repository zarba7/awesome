package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapGRPCLogger struct {
	lg    *zap.Logger
	sugar *zap.SugaredLogger
}

func (zl *zapGRPCLogger) Info(args ...interface{}) {
	if !zl.lg.Core().Enabled(zapcore.DebugLevel) {
		return
	}
	zl.sugar.Info(args...)
}

func (zl *zapGRPCLogger) Infoln(args ...interface{}) {
	if !zl.lg.Core().Enabled(zapcore.DebugLevel) {
		return
	}
	zl.sugar.Info(args...)
}

func (zl *zapGRPCLogger) Infof(format string, args ...interface{}) {
	if !zl.lg.Core().Enabled(zapcore.DebugLevel) {
		return
	}
	zl.sugar.Infof(format, args...)
}

func (zl *zapGRPCLogger) Warning(args ...interface{}) {
	zl.sugar.Warn(args...)
}

func (zl *zapGRPCLogger) Warningln(args ...interface{}) {
	zl.sugar.Warn(args...)
}

func (zl *zapGRPCLogger) Warningf(format string, args ...interface{}) {
	zl.sugar.Warnf(format, args...)
}

func (zl *zapGRPCLogger) Error(args ...interface{}) {
	zl.sugar.Error(args...)
}

func (zl *zapGRPCLogger) Errorln(args ...interface{}) {
	zl.sugar.Error(args...)
}

func (zl *zapGRPCLogger) Errorf(format string, args ...interface{}) {
	zl.sugar.Errorf(format, args...)
}

func (zl *zapGRPCLogger) Fatal(args ...interface{}) {
	zl.sugar.Fatal(args...)
}

func (zl *zapGRPCLogger) Fatalln(args ...interface{}) {
	zl.sugar.Fatal(args...)
}

func (zl *zapGRPCLogger) Fatalf(format string, args ...interface{}) {
	zl.sugar.Fatalf(format, args...)
}

func (zl *zapGRPCLogger) V(l int) bool {
	// infoLog == 0
	if l <= 0 { // debug level, then we ignore info level in gRPC
		return !zl.lg.Core().Enabled(zapcore.DebugLevel)
	}
	return true
}
