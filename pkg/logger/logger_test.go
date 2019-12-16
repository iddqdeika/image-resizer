package logger

import (
	"testing"
)

func TestConsoleLogger(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil{
			t.Errorf("must not panic")
		}
	}()

	log1 := ConsoleLogger(4)
	log1.Errorf("", nil)
	log1.Error(nil)
	log1.Infof("", nil)
	log1.Info(nil)

	log2 := ConsoleLogger(-4)
	log2.Errorf("", nil)
	log2.Error(nil)
	log2.Infof("", nil)
	log2.Info(nil)
}
