package easylog

import "strings"

type LogLevel int

const (
	DEBUG LogLevel = iota // 0
	INFO                  // 1
	WARN                  // 2
	ERROR                 // 3
	FATAL                 // 4
)

var levelFmt map[LogLevel]string

func init() {
	levelFmt = make(map[LogLevel]string)

	levelFmt[DEBUG] = "Debug"
	levelFmt[INFO] = "Info"
	levelFmt[WARN] = "Warn"
	levelFmt[ERROR] = "Error"
	levelFmt[FATAL] = "Fatal"
}

func GetLevelStr(level LogLevel) string {
	return levelFmt[level]
}

func GetLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}
