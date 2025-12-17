package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func NewLogger(logLevel string) (*zap.Logger, error) {
	var level zapcore.Level
	var encoding string
	var encodeLevel zapcore.LevelEncoder

	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
		encoding = "console"
		encodeLevel = zapcore.CapitalColorLevelEncoder
	case "info":
		level = zapcore.InfoLevel
		encoding = "json"
		encodeLevel = zapcore.LowercaseLevelEncoder
	default:
		level = zapcore.InfoLevel
		encoding = "json"
		encodeLevel = zapcore.LowercaseLevelEncoder
	}

	config := zap.Config{
		Encoding:         encoding,
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     customTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("error building zap logger: %w", err)
	}

	logger.Info("Logger initialized",
		zap.String("level", level.String()),
		zap.String("encoding", config.Encoding))

	return logger, nil
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("02-01-2006 15:04:05"))
}
