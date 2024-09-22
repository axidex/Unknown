package logger

type Logger interface {
	Info(keyAndValues ...interface{})
	Warn(keyAndValues ...interface{})
	Error(keyAndValues ...interface{})
	Fatal(keyAndValues ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type ConfigLogger struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}
