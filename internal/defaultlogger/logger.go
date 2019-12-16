package defaultlogger

import "fmt"

var Logger = &defLogger{}

type defLogger struct {

}

func (l *defLogger) Error(args ... interface{}) {
	fmt.Print(args...)
}

func (l *defLogger) Errorf(format string, args ... interface{}) {
	fmt.Printf(format, args)
}

func (l *defLogger) Info(args ... interface{}) {
	fmt.Print(args...)
}

func (l *defLogger) Infof(format string, args ... interface{}) {
	fmt.Printf(format, args)
}

