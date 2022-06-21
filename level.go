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

	levelFmt[DEBUG] = "Debug"
	levelFmt[INFO] = "Info"
	levelFmt[WARN] = "Warn"
	levelFmt[ERROR] = "Error"
	levelFmt[FATAL] = "Fatal"
}
