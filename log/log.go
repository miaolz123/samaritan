package log

import (
	"fmt"
	"time"
)

// Logger struct
type Logger struct {
	Name string
}

// New : get a logger
func New(name string) Logger {
	return Logger{Name: name}
}

// Do : log a message
func (logger Logger) Do(action string, price, amount float64, msgs ...interface{}) {
	now := time.Now().Format("01-02 15:04:05")
	msg := ""
	for _, m := range msgs {
		msg += fmt.Sprintf("%+v", m)
	}
	fmt.Printf("%s | %10s | %6s | %8.2f | %8.3f | %s\n", now, logger.Name, action, price, amount, msg)
}
