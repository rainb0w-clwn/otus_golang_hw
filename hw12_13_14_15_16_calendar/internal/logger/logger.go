package logger

import (
	"io"
	"log"
)

type Level int

const (
	Debug Level = iota
	Info
	Warning
	Error
)

type Logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warning(string, ...interface{})
	Error(string, ...interface{})
}
type SLogger struct {
	Level  Level
	logger *log.Logger
}

func New(level Level, writer io.Writer) Logger {
	return &SLogger{
		Level:  level,
		logger: log.New(writer, "", 0),
	}
}

func (l *SLogger) Debug(msg string, args ...interface{}) {
	if l.Level > Debug {
		return
	}
	l.log("DEBUG: ", msg, args...)
}

func (l *SLogger) Info(msg string, args ...interface{}) {
	if l.Level > Info {
		return
	}
	l.log("INFO: ", msg, args...)
}

func (l *SLogger) Warning(msg string, args ...interface{}) {
	if l.Level > Warning {
		return
	}
	l.log("WARN: ", msg, args...)
}

func (l *SLogger) Error(msg string, args ...interface{}) {
	if l.Level > Error {
		return
	}
	l.log("ERROR: ", msg, args...)
}

func (l *SLogger) log(prefix, msg string, args ...interface{}) {
	l.logger.SetPrefix(prefix)
	l.logger.Printf(msg+"\n", args...)
}
