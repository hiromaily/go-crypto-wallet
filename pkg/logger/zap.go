package logger

import (
	"io"
	"os"

	"github.com/yudai/pp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/hiromaily/go-bitcoin/pkg/config"
)

func NewLoggerWithWriter(w io.Writer, enab zapcore.LevelEnabler) *zap.Logger {
	pp.ColoringEnabled = false

	writer := zapcore.AddSync(w)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	cores := []zapcore.Core{
		zapcore.NewCore(jsonEncoder, writer, enab),
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
	return NewLoggerWithWriter(os.Stdout, getLogLevel(conf.Level)).Named(conf.Service)
}
