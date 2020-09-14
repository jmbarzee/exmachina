package system

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
)

const Version = "0.1.0"

// log is where normal & debugging messages are dumped to
var logger *log.Logger
var closeFile func() error

func Setup(logFilePath string) error {
	if logger != nil {
		return errors.New("system.Setup has already been called")
	}

	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	closeFile := logFile.Close
	logger = log.New(logFile, "", log.LstdFlags)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if err := closeFile(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
	return nil
}

func Panic(err error) {
	if closeFile != nil {
		if err := closeFile(); err != nil {
			panic(err)
		}
	}
	panic(err)
}

func Logf(fmt string, args ...interface{}) {
	if logger == nil {
		panic(errors.New("system.Setup has not been called"))
	}
	logger.Printf(fmt, args...)
}

func LogRoutinef(routineName, fmts string, args ...interface{}) {
	if logger == nil {
		panic(errors.New("system.Setup has not been called"))
	}
	prefix := fmt.Sprintf("Routine [%s]: ", routineName)
	logger.Printf(prefix+fmts, args...)
}

func LogRPCf(rpcName, fmts string, args ...interface{}) {
	if logger == nil {
		panic(errors.New("system.Setup has not been called"))
	}
	prefix := fmt.Sprintf("RPC-%s	: ", rpcName)
	logger.Printf(prefix+fmts, args...)
}

func Errorf(fmt string, args ...interface{}) {
	if logger == nil {
		panic(errors.New("system.Setup has not been called"))
	}
	logger.Printf("Error: "+fmt, args...)
}
