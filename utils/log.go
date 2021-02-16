package utils

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	logInfo = iota
	logWarn
	logError
)

type Log struct {
	path string
	mtx  sync.Mutex
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) SetPath(path string) {
	l.path = path
}

func (l *Log) LogMessage(module string, msg string, error string, logType int) error {
	var outMsg string
	var typ string
	dt := time.Now()

	var log, err = os.OpenFile(l.path+dt.Format("20060201")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println("Fail to open log file!")
		return err
	}

	switch logType {
	case logInfo:
		typ = "INFO"
		break

	case logWarn:
		typ = "WARN"
		break

	case logError:
		typ = "ERROR"
		break
	}

	outMsg = "[" + dt.Format("15:04:05") + "][" + module + "][" + typ + "] "
	outMsg += msg

	if logType == logError {
		outMsg += ": " + error
	}

	fmt.Println(outMsg)

	var _, errWr = log.WriteString(outMsg + "\n")
	if errWr != nil {
		fmt.Println("Fail to write to logfile!")
		log.Close()
		return errWr
	}

	log.Close()
	return nil
}

func (l *Log) Info(module string, msg string) error {
	var err error

	l.mtx.Lock()
	err = l.LogMessage(module, msg, "", logInfo)
	l.mtx.Unlock()

	return err
}

func (l *Log) Error(module string, msg string, err string) error {
	var er error

	l.mtx.Lock()
	er = l.LogMessage(module, msg, err, logError)
	l.mtx.Unlock()

	return er
}
