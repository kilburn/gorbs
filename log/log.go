package log

import (
	"fmt"
	"os"
	"sync"
)

const (
	// These levels are somewhat custom
	FATAL = iota
	ALERT
	ERROR
	WARN
	INFO
	DEBUG
)

// A Level represents the severity of a log message
type Level int

// A Logger represents an active logging object that generates lines of output
// to a number of backends. Multiple goroutines can use the logger at the same
// time.
type Logger struct {
	mu       sync.Mutex
	backends map[string]Backend
}

// NewLogger creates and returns a new logger with no associated backends.
func NewLogger() *Logger {
	return &Logger{backends: make(map[string]Backend)}
}

// write sends a message to all backends
func (l *Logger) write(level Level, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, backend := range l.backends {
		if backend.Level() >= level {
			backend.Write(level, message)
		}
	}
}

func (l *Logger) Log(level Level, v ...interface{}) {
	l.write(level, fmt.Sprint(v...))
}

func (l *Logger) Logf(level Level, format string, v ...interface{}) {
	l.write(level, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.Log(DEBUG, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logf(DEBUG, format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.Log(INFO, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logf(INFO, format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Log(WARN, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logf(WARN, format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Log(ERROR, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logf(ERROR, format, v...)
}

func (l *Logger) Alert(v ...interface{}) {
	l.Log(ALERT, v...)
	panic(fmt.Sprint(v...))
}

func (l *Logger) Alertf(format string, v ...interface{}) {
	l.Logf(ALERT, format, v...)
	panic(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Log(FATAL, v...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logf(FATAL, format, v...)
	os.Exit(1)
}

func (l *Logger) GetBackend(tag string) Backend {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.backends[tag]
}

func (l *Logger) AddBackend(tag string, backend Backend) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.backends[tag] = backend
}

func (l *Logger) RemoveBackend(tag string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.backends, tag)
}
