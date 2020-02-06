package easylog

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

	levelFmt[DEBUG] = "LV 0"
	levelFmt[INFO] = "LV 1"
	levelFmt[WARN] = "LV 2"
	levelFmt[ERROR] = "LV 3"
	levelFmt[FATAL] = "LV 4"
}
