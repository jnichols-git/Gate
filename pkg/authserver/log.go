package authserver

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var LogFile string = "./dat/log/authserver.log"
var logStream *os.File

func OpenLog() error {
	if logStream != nil {
		return errors.New("Attempted to OpenLog() before CloseLog() for previous open")
	}
	if file, err := os.OpenFile(LogFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend); err != nil {
		return err
	} else {
		logStream = file
	}
	return nil
}

func Log(format string, args ...interface{}) error {
	if logStream == nil {
		return errors.New("Attempted to Log() before OpenLog()")
	}
	currTime := time.Now().UTC().Format("06-01-02T15:04:05.999")
	logMsg := fmt.Sprintf("%s: %s", currTime, fmt.Sprintf(format, args...))
	_, err := logStream.Write([]byte(logMsg))
	return err
}

func CloseLog() {

}
