package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/hiromaily/go-bitcoin/pkg/config"
)

type LogEnv string

const (
	LogDev    LogEnv = "dev"
	LogProd   LogEnv = "prod"
	LogCustom LogEnv = "custom"
)

func (e LogEnv) String() string {
	return string(e)
}

func NewLoggerWithWriter(w io.Writer, lv zapcore.LevelEnabler, env LogEnv) *zap.Logger {
	//pp.ColoringEnabled = false

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
	//case LogProd:
	default:
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	cores := []zapcore.Core{
		zapcore.NewCore(jsonEncoder, writer, lv),
	}

	logger := zap.New(zapcore.NewTee(cores...),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

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

func NewZapLogger(conf *config.Logger) *zap.Logger {
	return NewLoggerWithWriter(os.Stdout, getLogLevel(conf.Level), LogEnv(conf.Env)).Named(conf.Service)
}
