package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var systemLog *zap.Logger
var requestLog *zap.Logger

func init() {
	{
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "system.log",
			MaxSize:    1 * 1024, // 1GB
			MaxBackups: 1,
			MaxAge:     7,
			Compress:   true,
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.InfoLevel,
		)
		systemLog = zap.New(core)
	}
	{
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "request.log",
			MaxSize:    1 * 1024, // 1GB
			MaxBackups: 1,
			MaxAge:     7,
			Compress:   true,
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.InfoLevel,
		)
		requestLog = zap.New(core)
	}
}

func SystemLog() *zap.Logger {
	return systemLog
}

func RequestLog() *zap.Logger {
	return requestLog
}
