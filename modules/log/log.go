package log

import (
	"log"
	"fmt"
	"time"
)

// Info records information to log.
func Info(format string, v ...interface{}) {
	now := time.Now()
	fmt.Printf("%04d/%02d/%02d %02d:%02d:%02d ", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	fmt.Printf(format + "\n", v ...)
}

// Error records error to log.
func Error(format string, v ...interface{}) {
	log.Printf(format + "\n", v ...)
}

// Fatal records error to log and exit process.
func Fatal(format string, v ...interface{}) {
	log.Fatalf(format + "\n", v ...)
}