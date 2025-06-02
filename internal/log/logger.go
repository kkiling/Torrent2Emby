package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	DebugLevel Level = Level(zapcore.DebugLevel)
	InfoLevel  Level = Level(zapcore.InfoLevel)
	WarnLevel  Level = Level(zapcore.WarnLevel)
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	With(args ...interface{}) Logger
	WithCtx(ctx context.Context, args ...interface{}) Logger
	Infow(string, ...interface{})
	Named(string) Logger
}

func NewLogger(level Level) Logger {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(level)),
		DisableCaller:    false,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			MessageKey:     "msg",
			LevelKey:       "level",
			NameKey:        "log",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			CallerKey:      "prepare",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
	}
	l, _ := cfg.Build()
	sugaredLogger := l.Sugar()
	defer sugaredLogger.Sync()
	return &logger{
		SugaredLogger: sugaredLogger,
	}
}

type logger struct {
	*zap.SugaredLogger
}

func (l *logger) With(args ...interface{}) Logger {
	return &logger{
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}

func (l *logger) WithCtx(ctx context.Context, args ...interface{}) Logger {
	return &logger{
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}

func (l *logger) Named(name string) Logger {
	return &logger{
		SugaredLogger: l.SugaredLogger.Named(name),
	}
}
