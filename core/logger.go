package core

import (
	"io"
	"log"
	"strings"

	"github.com/google/uuid"
)

type Logger struct {
	DebugMode bool
	*log.Logger
}

func NewLogger(out io.Writer, prefix string, flag int, debug bool) *Logger {
	return &Logger{DebugMode: debug, Logger: log.New(out, prefix, flag)}
}

func getLogID() string {
	logID := strings.Replace(uuid.New().String(), "-", "", -1)
	return logID
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

func (l *Logger) Print(v ...any) {
	l.Logger.Print(v...)
}

func (l *Logger) Printf(format string, v ...any) {
	l.Logger.Printf(format, v...)
}
