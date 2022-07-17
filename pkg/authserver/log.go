package authserver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var LogFile string = "./dat/log/authserver.log"
var openTime string
var logStream *os.File

// OpenLog opens a log file to target with Log
func OpenLog() error {
	if logStream != nil {
		return errors.New("Attempted to OpenLog() before CloseLog() for previous open")
	}
	if file, err := os.OpenFile(LogFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755); err != nil {
		return err
	} else {
		openTime = time.Now().UTC().Format("06-01-02T15:00")
		logStream = file
	}
	return nil
}

// Log writes a format string to the current open log file
func Log(format string, args ...interface{}) error {
	if logStream == nil {
		return errors.New("Attempted to Log() before OpenLog()")
	}
	currTime := time.Now().UTC().Format("06-01-02T15:04:05.999")
	logMsg := fmt.Sprintf("%s: %s\n", currTime, fmt.Sprintf(format, args...))
	_, err := logStream.Write([]byte(logMsg))
	return err
}

// CloseLog closes the current log file and rewrites it to a separate file
func CloseLog() {
	logStream.Close()
	fullLog, _ := ioutil.ReadFile(LogFile)
	permLogFile := fmt.Sprintf("./dat/log/authserver_%s.log", openTime)
	ioutil.WriteFile(permLogFile, fullLog, 0755)
}
