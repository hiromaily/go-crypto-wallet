package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
func Initialize() error {
	return initZapLoggers()
}

func initZapLoggers() error {

	//FIXME
	//var encoder zapcore.Encoder
	//switch env {
	//case enum.EnvDev:
	//	encoderCfg := zap.NewDevelopmentEncoderConfig()
	//	encoder = zapcore.NewConsoleEncoder(encoderCfg)
	//case enum.EnvProd:
	//	encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//default:
	//	return errors.Errorf("type should be set by [dev,prod] but %s is set", env)
	//}
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	//DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	//InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	//WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	//ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	//DPanicLevel
	// PanicLevel logs a message, then panics.
	//PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	//FatalLevel

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

// Fatal uses fmt.Sprintf to log a templated message.
func Fatal(args ...interface{}) {
	sugaredLogger.Fatal(args...)
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

// Fatalf uses fmt.Sprintf to log a templated message.
func Fatalf(format string, args ...interface{}) {
	sugaredLogger.Fatalf(format, args...)
}

// ErrorWithStack errorのstack traceを表示する. github.com/pkg/errorsで作成されたものしかstack表示できないので注意
// TODO:
func ErrorWithStack(error error) {
	noStackSugaredLogger.Errorf("error with stack trace: %+v", error)
}
