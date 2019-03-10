package log

import (
	"log"
)

// Error records error log.
func Error(format string, v ...interface{}) {
	log.Printf(format, v)
}

// Fatal records error log and exit process.
func Fatal(format string, v ...interface{}) {
	log.Fatalf(format, v)
}