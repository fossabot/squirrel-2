package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	// Log is the logger for normal use
	Log *log.Logger
	// Error is the Logger for errors
	Error *log.Logger
)

const errLogName = "error.log"

// Init creates logger instance to loggers
func Init() {
	initLogger()
	initErrorLogger()
}

func initLogger() {
	flag := log.Ldate | log.Ltime
	Log = log.New(os.Stdout, "", flag)
}

func initErrorLogger() {
	f, err := os.OpenFile(errLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	errHandler := io.MultiWriter(os.Stderr, f)
	flag := log.Ldate | log.Ltime | log.Lshortfile
	Error = log.New(errHandler, "", flag)
}

// UpdatePrefix Sets new prefix
func UpdatePrefix(prefix string) {
	if prefix != "" {
		prefix = fmt.Sprintf("[%s] ", prefix)
	}
	Log.SetPrefix(prefix)
	Error.SetPrefix(prefix)
}

// Printf is the alias for Log.Printf
func Printf(format string, v ...interface{}) {
	Log.Printf(format, v...)
}

// Println is the alias for Log.Printf
func Println(v ...interface{}) {
	Log.Println(v...)
}
