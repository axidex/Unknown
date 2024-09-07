package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	sugar  *zap.SugaredLogger
	logger *zap.Logger
}

func CreateNewZapLogger(loggerConfig ConfigLogger) (Logger, error) {
	level, err := zapcore.ParseLevel(loggerConfig.Level)
	if err != nil {
		return nil, err
	}

	atomicLevel := zap.NewAtomicLevelAt(level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // Capitalize the log level names
		EncodeTime:     zapcore.RFC3339TimeEncoder,     // RFC3339 UTC timestamp format
		EncodeDuration: zapcore.SecondsDurationEncoder, // Duration in seconds
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Short caller (file and line)
	}

	zapConfig := zap.Config{
		Level:            atomicLevel,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	return &ZapLogger{
		sugar:  sugar,
		logger: logger,
	}, nil
}

func (l *ZapLogger) Info(keyAndValues ...interface{}) {
	l.sugar.Info(keyAndValues...)
}

func (l *ZapLogger) Warn(keyAndValues ...interface{}) {
	l.sugar.Warn(keyAndValues...)
}

func (l *ZapLogger) Error(keyAndValues ...interface{}) {
	l.sugar.Error(keyAndValues...)
}

func (l *ZapLogger) Fatal(keyAndValues ...interface{}) {
	l.sugar.Fatal(keyAndValues...)
}

func (l *ZapLogger) Infof(msg string, args ...interface{}) {
	l.sugar.Infof(msg, args...)
}

func (l *ZapLogger) Warnf(msg string, args ...interface{}) {
	l.sugar.Warnf(msg, args...)
}

func (l *ZapLogger) Errorf(msg string, args ...interface{}) {
	l.sugar.Errorf(msg, args...)
}

func (l *ZapLogger) Fatalf(msg string, args ...interface{}) {
	l.sugar.Fatalf(msg, args...)
}
