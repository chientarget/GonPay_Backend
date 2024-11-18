package logger

import (
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

type logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
}

var (
	instance Logger
	once     sync.Once
)

func NewLogger() Logger {
	once.Do(func() {
		instance = &logger{
			debug:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
			info:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			warning: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
			error:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		}
	})
	return instance
}

func (l *logger) Debug(msg string, keyvals ...interface{}) {
	l.debug.Printf(msg, keyvals...)
}

func (l *logger) Info(msg string, keyvals ...interface{}) {
	l.info.Printf(msg, keyvals...)
}

func (l *logger) Warn(msg string, keyvals ...interface{}) {
	l.warning.Printf(msg, keyvals...)
}

func (l *logger) Error(msg string, keyvals ...interface{}) {
	l.error.Printf(msg, keyvals...)
}
