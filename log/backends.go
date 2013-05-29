package log

import (
	"log"
	"log/syslog"
)

// A Backend knows how to output a message to some backend, such as stdout,
// syslog or a log file.
type Backend interface {
	Write(level Level, message string)
	Leveled
}

type Leveled interface {
	Level() Level
	SetLevel(level Level)
}

type LevelTracker struct {
	level Level
}

func (l *LevelTracker) Level() Level {
	return l.level
}

func (l *LevelTracker) SetLevel(level Level) {
	if level > DEBUG {
		level = DEBUG
	}
	l.level = level
}

// A NullBackend discards all messages instead of writing them
type NullBackend struct {
	*LevelTracker
}

func NewNullBackend(level Level) *NullBackend {
	backend := &NullBackend{&LevelTracker{}}
	backend.SetLevel(level)
	return backend
}

// Write doesn't do anything at all
func (w *NullBackend) Write(level Level, message string) {}

// An StdOutBackend writes log messages to the standard output (through the
// standard log package)
type StdBackend struct {
	*LevelTracker
	logger *log.Logger
}

func NewStdBackend(level Level, logger *log.Logger) *StdBackend {
	backend := &StdBackend{&LevelTracker{}, logger}
	backend.SetLevel(level)
	return backend
}

// Write sends a message to stdout.
func (w *StdBackend) Write(level Level, message string) {
	w.logger.Print(message)
}

// A SyslogBackend writes log messages to syslog (through the standard syslog
// package)
type SyslogBackend struct {
	*LevelTracker
	logger *syslog.Writer
}

func NewSyslogBackend(level Level, writer *syslog.Writer) *SyslogBackend {
	backend := &SyslogBackend{&LevelTracker{}, writer}
	backend.SetLevel(level)
	return backend
}

// Write sends the message through one of the syslog.Writer functions,
// depending on the message level
func (w *SyslogBackend) Write(level Level, message string) {
	switch level {
	case FATAL:
		w.logger.Crit(message)
	case ALERT:
		w.logger.Alert(message)
	case ERROR:
		w.logger.Err(message)
	case WARN:
		w.logger.Warning(message)
	case INFO:
		w.logger.Info(message)
	case DEBUG:
		w.logger.Debug(message)
	default:
		w.logger.Notice(message)
	}
}
