package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var Logger *zap.Logger

var once = sync.Once{}

func init() {
	once.Do(func() {
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		defaultEncoder := zapcore.NewJSONEncoder(cfg)

		// todo add to config
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./logs/log.json",
			MaxSize:    1,
			MaxAge:     7,
			MaxBackups: 7,
			LocalTime:  false,
			Compress:   true,
		})

		stdoutWriter := zapcore.AddSync(os.Stdout)
		defaultLogLevel := zap.InfoLevel
		core := zapcore.NewTee(
			zapcore.NewCore(defaultEncoder, writer, defaultLogLevel),
			zapcore.NewCore(defaultEncoder, stdoutWriter, defaultLogLevel),
		)

		Logger = zap.New(core)
	})
}
