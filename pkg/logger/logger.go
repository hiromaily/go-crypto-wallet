package logger

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/status"
)

var (
	zapLogger            *zap.Logger
	sugaredLogger        *zap.SugaredLogger
	noStackSugaredLogger *zap.SugaredLogger
)

type nostackLevelEnabler struct{}

func (le *nostackLevelEnabler) Enabled(zapcore.Level) bool {
	return false
}

// Initialize 設定を読み込みLoggerを初期化する。mainクラスなど起動時に呼ばれる想定
func Initialize(env string) error {
	return initZapLoggers(env)
}

func initZapLoggers(env string) error {

	var encoder zapcore.Encoder
	switch env {
	case "dev":
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	case "prod":
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	default:
		return errors.Errorf("type should be set by [dev,prod] but %s is set", env)
	}

	error := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl >= zapcore.ErrorLevel })
	info := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl < zapcore.ErrorLevel })

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), error),
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), info))

	zl := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel))

	nostack := zl.WithOptions(zap.AddStacktrace(&nostackLevelEnabler{})).Sugar()

	zapLogger = zl
	sugaredLogger = zl.Sugar()
	noStackSugaredLogger = nostack

	return nil
}

// Debug fmt.Sprintf to log a templated message.
func Debug(args ...interface{}) {
	sugaredLogger.Debug(args...)
}

// Info uses fmt.Sprintf to log a templated message.
func Info(args ...interface{}) {
	sugaredLogger.Info(args...)
}

// Warn uses fmt.Sprintf to log a templated message.
func Warn(args ...interface{}) {
	sugaredLogger.Warn(args...)
}

// Error uses fmt.Sprintf to log a templated message.
func Error(args ...interface{}) {
	sugaredLogger.Error(args...)
}

// Debugf fmt.Sprintf to log a templated message.
func Debugf(format string, args ...interface{}) {
	sugaredLogger.Debugf(format, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(format string, args ...interface{}) {
	sugaredLogger.Infof(format, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(format string, args ...interface{}) {
	sugaredLogger.Warnf(format, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(format string, args ...interface{}) {
	sugaredLogger.Errorf(format, args...)
}

// ErrorWithStack errorのstack traceを表示する. github.com/pkg/errorsで作成されたものしかstack表示できないので注意
func ErrorWithStack(error error) {
	st := status.Convert(error)
	noStackSugaredLogger.Errorf("error with stack trace: %+v, details = %+v", error, st.Details())
}
