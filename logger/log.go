package logger

import (
	"log"
	"os"
)

var Log *log.Logger

func init() {
	// open logger
	Log = log.New(os.Stdout, "" /*log.Lshortfile|*/, log.Ltime)
}
