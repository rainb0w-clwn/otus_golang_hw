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
	Debug(string)
	Info(string)
	Warning(string)
	Error(string)
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

func (l *SLogger) Debug(msg string) {
	if l.Level > Debug {
		return
	}
	l.log("DEBUG: ", msg)
}

func (l *SLogger) Info(msg string) {
	if l.Level > Info {
		return
	}
	l.log("INFO: ", msg)
}

func (l *SLogger) Warning(msg string) {
	if l.Level > Warning {
		return
	}
	l.log("WARN: ", msg)
}

func (l *SLogger) Error(msg string) {
	if l.Level > Error {
		return
	}
	l.log("ERROR: ", msg)
}

func (l *SLogger) log(prefix, msg string) {
	l.logger.SetPrefix(prefix)
	l.logger.Println(msg)
}
