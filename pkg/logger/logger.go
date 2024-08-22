package logger

type Level int32

const (
	SILENT Level = iota
	PANIC
	ERROR
	WARN
	INFO
	DEBUG
)

type Logger interface {
	Log(lvl Level, msg string)
}

func (l Level) String() string {
	return [...]string{"SILENT", "PANIC", "ERROR", "WARN", "INFO", "DEBUG", "VERBOSE"}[l]
}

var logger Logger

func SetLogger(lg Logger) {
	logger = lg
}

func Log(lvl Level, msg string) {
	if logger != nil {
		logger.Log(lvl, msg)
	}
}
