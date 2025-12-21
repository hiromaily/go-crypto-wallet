package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogEnv dev/prod/custom
type LogEnv string

// log_type
const (
	LogDev    LogEnv = "dev"
	LogProd   LogEnv = "prod"
	LogCustom LogEnv = "custom"
)

// String converter
func (e LogEnv) String() string {
	return string(e)
}

// NewLoggerWithWriter returns *zap.Logger
func NewLoggerWithWriter(w io.Writer, env LogEnv, lv zapcore.LevelEnabler, isStackTrace bool) *zap.Logger {
	zap.NewExample()

	writer := zapcore.AddSync(w)

	var encoderCfg zapcore.EncoderConfig
	switch env {
	case LogDev:
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	case LogCustom:
		encoderCfg = zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "lv",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
	case LogProd:
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	default:
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	cores := []zapcore.Core{
		zapcore.NewCore(jsonEncoder, writer, lv),
	}
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}
	if isStackTrace {
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}
	logger := zap.New(zapcore.NewTee(cores...), options...)

	return logger
}

func getLogLevel(level string) zapcore.LevelEnabler {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// NewZapLogger returns *zap.Logger
func NewZapLogger(env, level, service string, isLogger bool) *zap.Logger {
	return NewLoggerWithWriter(
		os.Stdout,
		LogEnv(env),
		getLogLevel(level),
		isLogger).Named(service)
}
