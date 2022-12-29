package easylog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultLevel   LogLevel = DEBUG
	defaultMaxSize          = 100 << (10 * 2) // MB
	defaultDir              = ""
)

var (
	logInit  bool = false
	level    LogLevel
	logger   *log.Logger
	fp       *os.File
	maxSize  int64
	dir      string
	fileName string
	mutex    sync.Mutex
)

func init() {
	level = defaultLevel
	maxSize = defaultMaxSize
	dir = defaultDir
	fileName = fmt.Sprintf("%s.log", filepath.Base(os.Args[0]))
}

func Logger() *log.Logger {
	return logger
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Init(options ...func()) error {
	if logInit {
		return nil
	}

	for _, f := range options {
		f()
	}
	if err := checkFile(); err != nil {
		return err
	}

	if logger == nil {
		return errors.New("nil logger")
	}

	logInit = true
	return nil
}

func Debug(a ...any) {
	if level <= DEBUG {
		check(checkFile())
		msg := fmt.Sprint(a...)
		fmtMsg := format(DEBUG, msg)
		fmt.Println(fmtMsg)
		logger.Println(fmtMsg)
	}
}

func Info(a ...any) {
	if level <= INFO {
		check(checkFile())
		msg := fmt.Sprint(a...)
		fmtMsg := format(INFO, msg)
		fmt.Println(fmtMsg)
		logger.Println(fmtMsg)
	}
}

func Warn(a ...any) {
	if level <= WARN {
		check(checkFile())
		msg := fmt.Sprint(a...)
		fmtMsg := format(WARN, msg)
		fmt.Println(fmtMsg)
		logger.Println(fmtMsg)
	}
}

func Error(a ...any) {
	if level <= ERROR {
		check(checkFile())
		msg := fmt.Sprint(a...)
		fmtMsg := format(ERROR, msg)
		fmt.Println(fmtMsg)
		logger.Println(fmtMsg)
	}
}

func Fatal(a ...any) {
	if level <= FATAL {
		check(checkFile())
		msg := fmt.Sprint(a...)
		fmtMsg := format(FATAL, msg)
		fmt.Println(fmtMsg)
		logger.Println(fmtMsg)
	}
}

func checkFile() error {
	mutex.Lock()
	defer mutex.Unlock()

	if fp == nil {
		if err := openFile(dir, fileName); err != nil {
			return err
		}
	}

	if isFileMax(fp) {
		closeFile()
		if err := renameFile(); err != nil {
			return err
		}
		if err := openFile(dir, fileName); err != nil {
			return err
		}
		setNewLogger(fp)
	}

	if logger == nil {
		setNewLogger(fp)
	}
	return nil
}

func renameFile() error {
	old := filepath.Join(dir, fileName)
	new := fmt.Sprintf("%s.bak.%s", filepath.Join(dir, fileName), time.Now().Format("20060102150405"))

	if err := os.Rename(old, new); err != nil {
		return err
	}
	return nil
}

func setNewLogger(fp *os.File) {
	logger = log.New(fp, "", 0)
}

func closeFile() {
	if fp == nil {
		return
	}
	fp.Close()
	fp = nil
}

func isFileMax(fp *os.File) bool {
	info, err := fp.Stat()
	if err != nil {
		panic(err)
	}

	if info.Size() >= maxSize {
		return true
	}
	return false
}

func openFile(dir, name string) error {
	var err error

	fp, err = os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	return nil
}

func GetLevel() LogLevel {
	return level
}

func SetLevel(lv LogLevel) func() {
	return func() {
		level = lv
	}
}

func SetMaxSize(size int) func() {
	return func() {
		maxSize = int64(size) << (10 * 2) // MB to Byte
	}
}

func SetDir(newDir string) func() {
	return func() {
		dir = newDir
	}
}

func SetFileName(name string) func() {
	return func() {
		fileName = name
	}
}
