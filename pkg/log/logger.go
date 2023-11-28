package log

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerMode string

type Config struct {
	AppName string
	Mode    LoggerMode
	LokiURL string
}

func NewLogger(args Config) (*zap.Logger, error) {
	outputPath := []string{"stdout"}
	errOutputPath := []string{"stderr"}
	zapCfg := zap.Config{
		Encoding:         "json",
		OutputPaths:      outputPath,
		ErrorOutputPaths: errOutputPath,
	}

	switch args.Mode {
	case "prod":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		zapCfg.EncoderConfig = encoderProdConfig
	case "dev":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		zapCfg.EncoderConfig = encoderDevConfig
	default:
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		zapCfg.EncoderConfig = encoderDevConfig
	}

	loki := New(context.Background(), LokiConfig{
		Url:          args.LokiURL,
		BatchMaxSize: 1000,
		BatchMaxWait: 10 * time.Second,
		Labels:       map[string]string{"app": "grafanacloud-ngdlong91-logs"},
	})
	return loki.WithCreateLogger(zapCfg)

}

var encoderProdConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "severity",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    "function",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    encodeLevel(),
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.MillisDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeName:     ShortNameEncoder,
}

var encoderDevConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "severity",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    encodeLevel(),
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.MillisDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	//EncodeName:     ShortNameEncoder,
}

func ShortNameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	loggerNameArr := strings.Split(loggerName, "/")
	fmt.Println()
	enc.AppendString(loggerNameArr[len(loggerNameArr)-1])
}

func encodeLevel() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.PanicLevel:
			enc.AppendString("ALERT")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		}
	}
}

func LoggerForPackage(lgr *zap.Logger, packageName string) *zap.Logger {
	return lgr.With(zap.String("package", packageName))
}

func LoggerForFunction(lgr *zap.Logger, funcName string) *zap.Logger {
	return lgr.With(zap.String("function", funcName))
}

func LoggerForClass(lgr *zap.Logger, className string) *zap.Logger {
	return lgr.With(zap.String("class", className))
}
