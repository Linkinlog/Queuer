package logger

import (
	"io"

	"github.com/linkinlog/queuer/internal/db"
)

type LoggerParams struct {
	W   io.Writer
	DSN string
}

func New(l *LoggerParams) *Logger {
	return &Logger{
		Writer: l.W,
		dsn:    l.DSN,
	}
}

type Logger struct {
	io.Writer
	dsn string
}

func (l *Logger) Write(p []byte) (n int, err error) {
	n, err = l.Writer.Write(p)
	if err != nil {
		return 0, err
	}

	dbLogger, err := db.OpenLog(l.dsn)
	if err != nil {
		return 0, err
	}

	if err := dbLogger.WriteLog(p); err != nil {
		return 0, err
	}

	return n, nil
}
