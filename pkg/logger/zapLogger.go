package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"syscall"
	"time"
)

// ZapLogger Logger
type ZapLogger struct {
	sugarLogger *zap.SugaredLogger
	RawLogger   *zap.Logger
	cfg         ConfigLogger
}

func CreateNewZapLogger(config ConfigLogger) (Logger, error) {
	// file
	file, err := os.OpenFile(config.FilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	logger := &ZapLogger{cfg: config}
	logger.InitLogger(file)

	return logger, nil
}

// For mapping config logger to app logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *ZapLogger) getLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[l.cfg.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// CustomTimeEncoder Custom time encoder
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000-07:00"))
}

// CustomLevelEncoder Custom level encoder with uppercase
func CustomLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}

// InitLogger Init logger
func (l *ZapLogger) InitLogger(f *os.File) {
	logLevel := l.getLoggerLevel()

	var encoderCfg zapcore.EncoderConfig

	encoderCfg = zap.NewProductionEncoderConfig()

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "level"
	encoderCfg.CallerKey = "caller"
	encoderCfg.TimeKey = "time"
	encoderCfg.NameKey = "name"
	encoderCfg.MessageKey = "message"
	encoderCfg.EncodeTime = CustomTimeEncoder
	encoderCfg.EncodeLevel = CustomLevelEncoder

	encoder = zapcore.NewConsoleEncoder(encoderCfg)
	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), zap.NewAtomicLevelAt(logLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(logLevel)),
	)

	l.RawLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = l.RawLogger.Sugar()
	err := l.sugarLogger.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		l.sugarLogger.Error(err)
	}
}

// Logger methods

func (l *ZapLogger) Info(keysAndValues ...interface{}) {
	l.sugarLogger.Info(keysAndValues...)
}

func (l *ZapLogger) Warn(keysAndValues ...interface{}) {
	l.sugarLogger.Warn(keysAndValues...)
}

func (l *ZapLogger) Error(keysAndValues ...interface{}) {
	l.sugarLogger.Error(keysAndValues...)
}

func (l *ZapLogger) Fatal(keysAndValues ...interface{}) {
	l.sugarLogger.Fatal(keysAndValues...)
}

func (l *ZapLogger) Infof(msg string, args ...interface{}) {
	l.sugarLogger.Infof(msg, args...)
}
func (l *ZapLogger) Warnf(msg string, args ...interface{}) {
	l.sugarLogger.Warnf(msg, args...)
}
func (l *ZapLogger) Errorf(msg string, args ...interface{}) {
	l.sugarLogger.Errorf(msg, args...)
}
func (l *ZapLogger) Fatalf(msg string, args ...interface{}) {
	l.sugarLogger.Fatalf(msg, args...)
}
