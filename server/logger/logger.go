package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	file *os.File
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{file: file}, nil
}

func (l *Logger) Log(message string) {
	fmt.Fprintln(l.file, message)
}

func (l *Logger) Close() {
	l.file.Close()
}
