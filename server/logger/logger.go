package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	file      *os.File
	writeChan chan string
	done      chan struct{}
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	l := &Logger{
		file:      file,
		writeChan: make(chan string),
		done:      make(chan struct{}),
	}

	go l.writeLoop()

	return l, nil
}

func (l *Logger) Log(message string) {
	l.writeChan <- message
}

func (l *Logger) Close() {
	close(l.writeChan)
	<-l.done
}

func (l *Logger) writeLoop() {
	for message := range l.writeChan {
		fmt.Fprintln(l.file, message)
	}
	l.file.Close()
	close(l.done)
}
