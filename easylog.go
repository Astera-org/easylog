package easylog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	defaultLevel        = DEBUG
	defaultMaxSizeBytes = int64(100 << (10 * 2)) // 100 MiB
)

var (
	level    LogLevel = defaultLevel
	logger   *log.Logger
	fp       *os.File
	maxSize  int64 = defaultMaxSizeBytes
	filePath string
	mutex    sync.Mutex
	stdout   *os.File = os.Stdout
)

// logFilePath = "" means only log to stdout
func Init(logFilePath string) error {
	closeFile()
	filePath = logFilePath
	return checkFile(false)
}

func Logger() *log.Logger {
	return logger
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func logAt(l LogLevel, a ...any) {
	if l < level {
		return
	}
	check(checkFile(false))
	sprintArgs := []interface{}{levelFmt[l], ": "}
	sprintArgs = append(sprintArgs, a...)
	msg := fmt.Sprint(sprintArgs...)
	logger.Output(3, msg)
	if l == FATAL {
		panic(msg)
	}
}

func logAtf(l LogLevel, format string, a ...any) {
	if l < level {
		return
	}
	check(checkFile(false))
	msg := fmt.Sprintf("%s: %s", levelFmt[l], fmt.Sprintf(format, a...))
	logger.Output(3, msg)
	if l == FATAL {
		panic(msg)
	}
}

func Debug(a ...any) {
	logAt(DEBUG, a...)
}

func Info(a ...any) {
	logAt(INFO, a...)
}

func Warn(a ...any) {
	logAt(WARN, a...)
}

func Error(a ...any) {
	logAt(ERROR, a...)
}

func Fatal(a ...any) {
	logAt(FATAL, a...)
}

func Debugf(format string, a ...any) {
	logAtf(DEBUG, format, a...)
}

func Infof(format string, a ...any) {
	logAtf(INFO, format, a...)
}

func Warnf(format string, a ...any) {
	logAtf(WARN, format, a...)
}

func Errorf(format string, a ...any) {
	logAtf(ERROR, format, a...)
}

func Fatalf(format string, a ...any) {
	logAtf(FATAL, format, a...)
}

func checkFile(locked bool) error {
	if !locked {
		mutex.Lock()
		defer mutex.Unlock()
	}
	var err error

	if fp == nil {
		var writer io.Writer
		if filePath == "" {
			fp = stdout
			writer = stdout
		} else {
			if fp, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
				return err
			}
			writer = io.MultiWriter(stdout, fp)
		}
		logger = log.New(writer, "", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	}

	info, err := fp.Stat()
	if err != nil {
		return err
	}

	if info.Size() >= maxSize {
		closeFile()
		new := fmt.Sprintf("%s.bak.%s", filePath, time.Now().Format("20060102150405"))
		if err = os.Rename(filePath, new); err != nil {
			return err
		}
		return checkFile(true)
	}

	return nil
}

func closeFile() {
	if fp == nil {
		return
	}
	fp.Close()
	fp = nil
}

func GetLevel() LogLevel {
	return level
}

func SetLevel(lv LogLevel) {
	level = lv
}

func SetMaxSize(sizeMiB int) {
	if sizeMiB <= 0 {
		maxSize = defaultMaxSizeBytes
		Error("SetMaxSize: sizeMiB must be greater than 0, using default value: ", defaultMaxSizeBytes)
	} else {
		maxSize = int64(sizeMiB) << (10 * 2) // MiB to Byte
	}
}
