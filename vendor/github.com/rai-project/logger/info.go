package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// returns file and line two stack frames above its invocation
func Caller() (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return
}

// Decorate appends line, file and function context to the logger and returns a function to call before
// each log
func Decorate(logger *logrus.Entry) func() *logrus.Entry {
	return func() *logrus.Entry {
		if pc, f, line, ok := runtime.Caller(1); ok {
			fnName := runtime.FuncForPC(pc).Name()
			file := strings.Split(f, "mobilebid")[1]
			caller := fmt.Sprintf("%s:%v %s", file, line, fnName)

			return logrus.WithField("caller", caller)
		}
		return logger
	}
}

// Log appends line, file and function context to the logger
func Log() *logrus.Entry {
	if pc, f, line, ok := runtime.Caller(1); ok {
		fnName := runtime.FuncForPC(pc).Name()
		file := strings.Split(f, "mobilebid")[1]
		caller := fmt.Sprintf("%s:%v %s", file, line, fnName)

		return logrus.WithField("caller", caller)
	}
	return &logrus.Entry{}
}
