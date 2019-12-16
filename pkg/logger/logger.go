package logger

import (
	"fmt"
)

const(
	LogErrors LogLevel = iota
	LogAll
)

type LogLevel int

type Logger interface {
	Error(args... interface{})
	Errorf(format string, args... interface{})
	Info(args... interface{})
	Infof(format string, args... interface{})
}

func ConsoleLogger(level LogLevel) Logger{
	return &consoleLogger{
		l:	level,
	}
}

type consoleLogger struct {
	l		LogLevel
}

func (l *consoleLogger) Error(args ...interface{}) {
	if l.l>=LogErrors{
		fmt.Print(args...)
	}
}

func (l *consoleLogger) Errorf(format string, args ...interface{}) {
	if l.l>=LogErrors{
		fmt.Printf(format, args...)
	}
}

func (l *consoleLogger) Info(args ... interface{}) {
	if l.l>=LogAll{
		fmt.Print(args...)
	}
}

func (l *consoleLogger) Infof(format string, args ... interface{}) {
	if l.l>=LogAll{
		fmt.Printf(format, args...)
	}
}

