package utils

import (
	"log"
	"os"
)

const (
	INFO = iota
	DEBUG
	WARN
	ERROR
)

type Logger struct {
	level  int
	logger *log.Logger
}

func NewLogger(level int) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Debug(msg string, v ...interface{}) {
	if l.level <= DEBUG {
		l.logger.SetPrefix("DEBUG: ")
		l.logger.Printf(msg, v...)
	}
}

func (l *Logger) Info(msg string, v ...interface{}) {
	if l.level <= INFO {
		l.logger.SetPrefix("INFO: ")
		l.logger.Printf(msg, v...)
	}
}

func (l *Logger) Warn(msg string, v ...interface{}) {
	if l.level <= WARN {
		l.logger.SetPrefix("WARN: ")
		l.logger.Printf(msg, v...)
	}
}

func (l *Logger) Error(msg string, v ...interface{}) {
	if l.level <= ERROR {
		l.logger.SetPrefix("ERROR: ")
		l.logger.Printf(msg, v...)
	}
}
