package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// create a logger for error, warning, and info
var errorL = log.New(newLogWriter("[ERROR]"), "", 0)
var warnL = log.New(newLogWriter("[WARNING]"), "", 0)
var infL = log.New(newLogWriter("[INFO]"), "", 0)

// a function that logs an error
func Error(err error) {
	errorL.Println(err)
}

// a function that logs a warning
func Warn(warning string) {
	warnL.Println(warning)
}

// a function that logs info
func Info(info string) {
	infL.Println(info)
}

// a custom writer that formats the log messages
type logWriter struct {
	prefix string // the prefix to add before the message
}

// a function that creates a new logWriter with the given prefix
func newLogWriter(prefix string) *logWriter {
	return &logWriter{prefix: prefix}
}

func (writer *logWriter) Write(bytes []byte) (int, error) {
	time := time.Now().Format("15:04:05 02.01.2006") // this is for some reason the golang way to format time
	output := fmt.Sprintf("%s %s %s", time, writer.prefix, string(bytes))
	return os.Stderr.Write([]byte(output))
}
