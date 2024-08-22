package logger

import (
	"log"

	"github.com/fatih/color"
)

type Logcat struct {
	lvl Level
}

func NewLogcat(lvl Level) Logger {
	return &Logcat{
		lvl: lvl,
	}
}

func (l *Logcat) Log(lvl Level, msg string) {
	if lvl <= l.lvl && lvl != SILENT {
		c := l.getColor(lvl)
		log.Printf("%s%s %s\t%s\n", color.YellowString("Vending"), color.CyanString("Machine"), lvl, c(msg))
	}
}

func (l *Logcat) getColor(lvl Level) func(...interface{}) string {
	switch lvl {
	case PANIC:
		return color.New(color.FgWhite, color.BgRed).SprintFunc()
	case ERROR:
		return color.New(color.FgRed).SprintFunc()
	case WARN:
		return color.New(color.FgYellow).SprintFunc()
	case DEBUG:
		return color.New(color.FgGreen).SprintFunc()
	default:
		return color.New(color.FgWhite).SprintFunc()
	}
}
