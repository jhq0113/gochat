package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger ILog
)

type ILog func(lvl zapcore.Level, msg string, fields ...zap.Field)

func Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger(lvl, msg, fields...)
	}
}
