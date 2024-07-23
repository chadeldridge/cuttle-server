package core

import (
	"io"
	"log"
)

type Logger struct {
	DebugMode bool
	*log.Logger
}

func NewLogger(out io.Writer, prefix string, flag int, debug bool) *Logger {
	return &Logger{DebugMode: debug, Logger: log.New(out, prefix, flag)}
}

func (l *Logger) Debug(v ...any) {
	if l.DebugMode {
		l.Print(append([]any{"[DEBUG] "}, v...)...)
	}
}

func (l *Logger) Debugf(format string, v ...any) {
	if l.DebugMode {
		l.Printf("[DEBUG] "+format, v...)
	}
}
