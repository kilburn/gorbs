package log

import (
	"log"
	"os"
)

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger()
	stdLogger := log.New(os.Stderr, "", log.LstdFlags)
	defaultLogger.AddBackend("std", NewStdBackend(WARN, stdLogger))
}

func Log(level Level, v ...interface{}) {
	defaultLogger.Log(level, v...)
}

func Logf(level Level, format string, v ...interface{}) {
	defaultLogger.Logf(level, format, v...)
}

func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func Alert(v ...interface{}) {
	defaultLogger.Alert(v...)
}

func Alertf(format string, v ...interface{}) {
	defaultLogger.Alertf(format, v...)
}

// Aliases to allow drop-in compatibility with the standard log package
var Panic, Panicf = Alert, Alertf

func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

func GetBackend(tag string) Backend {
	return defaultLogger.GetBackend(tag)
}

func AddBackend(tag string, backend Backend) {
	defaultLogger.AddBackend(tag, backend)
}

func RemoveBackend(tag string) {
	defaultLogger.RemoveBackend(tag)
}
