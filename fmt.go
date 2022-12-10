package easylog

import (
	"time"
)

const (
	timeFmt = "01-02 15:04:05"
)

func format(lv LogLevel, msg string) string {
	return "" + time.Now().Format(timeFmt) + ", " + levelFmt[lv] + ": " + msg
}
