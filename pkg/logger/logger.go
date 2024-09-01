package logger

type Logger interface {
	Info(msg string, keyAndValues ...interface{})
	Warn(msg string, keyAndValues ...interface{})
	Error(msg string, keyAndValues ...interface{})
	Fatal(msg string, keyAndValues ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}
