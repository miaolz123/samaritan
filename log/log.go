package log

import (
	"fmt"
	"time"
)

// Logger : logger struct
type Logger struct {
	Name string
}

// New : get a logger
func New(name string) Logger {
	return Logger{Name: name}
}

// Do : log a message
func (logger Logger) Do(subject, action, message string, price, amount float64) {
	now := time.Now().Format("01-02 15:04:05")
	fmt.Printf("%s | %10s | %6s | %8.2f | %8.3f | %s\n", now, subject, action, price, amount, message)
}
