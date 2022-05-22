package log

import (
	"fmt"
	"git-svn-bridge/conf"
	"log"
	"os"
	"runtime/debug"
)

var logFile *os.File
var fileErrorLogger *log.Logger
var stdErrorLogger *log.Logger

func InitLogging() {
	config := conf.GetConfig()

	var err error
	if logFile, err = os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664); err != nil {
		panic(fmt.Errorf("could not open log file: %w", err))
	}

	fileErrorLogger = log.New(logFile, "", log.LstdFlags)
	stdErrorLogger = log.New(os.Stderr, "", log.LstdFlags)
}

func Fatalf(format string, v ...interface{}) {
	fileErrorLogger.Fatalf(format, v...)
}

func StdErrFatalf(format string, v ...interface{}) {
	stdErrorLogger.Printf(format, v...)
	fileErrorLogger.Fatalf(format, v...)
}

func StdErrPrintf(format string, v ...interface{}) {
	stdErrorLogger.Printf(format, v...)
	fileErrorLogger.Printf(format, v...)
}

func OnPanicf(err error) {
	if r := recover(); r != nil {
		stdErrorLogger.Println(err)
		Fatalf("FATAL: %s. Cause: %s\n%s", err, r, debug.Stack())
	}
}

func StdErrOnPanicf(err error) {
	if r := recover(); r != nil {
		StdErrFatalf("FATAL: %s. Cause: %s\n%s", err, r, debug.Stack())
	}
}

func CloseLog() {
	if logFile != nil {
		err := logFile.Sync()
		if err != nil {
			panic(fmt.Errorf("could not sync log file before close: %w", err))
		}

		err = logFile.Close()
		if err != nil {
			panic(fmt.Errorf("could not close log file: %w", err))
		}
	}
}
