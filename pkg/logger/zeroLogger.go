package logger

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

type ZeroLogger struct {
	logger zerolog.Logger
}

func CreateNewZeroLogger(loggerConfig ConfigLogger) (Logger, error) {
	level, err := zerolog.ParseLevel(loggerConfig.Level)
	if err != nil {
		return nil, err
	}
	zerolog.SetGlobalLevel(level)

	// file
	file, err := os.OpenFile(loggerConfig.FilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// console
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC822,
	}

	logger := zerolog.New(io.MultiWriter(consoleWriter, file)).With().Timestamp().Logger()

	return &ZeroLogger{logger: logger}, nil
}

func (l *ZeroLogger) Info(keyAndValues ...interface{}) {
	l.logger.Info().Fields(keyAndValues)
}

func (l *ZeroLogger) Warn(keyAndValues ...interface{}) {
	l.logger.Warn().Fields(keyAndValues)
}

func (l *ZeroLogger) Error(keyAndValues ...interface{}) {
	l.logger.Error().Fields(keyAndValues)
}

func (l *ZeroLogger) Fatal(keyAndValues ...interface{}) {
	l.logger.Fatal().Fields(keyAndValues)
}

func (l *ZeroLogger) Infof(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

func (l *ZeroLogger) Warnf(msg string, args ...interface{}) {
	l.logger.Warn().Msgf(msg, args...)
}

func (l *ZeroLogger) Errorf(msg string, args ...interface{}) {
	l.logger.Error().Msgf(msg, args...)
}

func (l *ZeroLogger) Fatalf(msg string, args ...interface{}) {
	l.logger.Fatal().Msgf(msg, args...)
}
