package ui

import (
	"log"
	"os"
	"sync"
)

var logger *log.Logger
var once sync.Once
var filename = "amnisiac.log"

// GetLoggerInstance returns the UI global logger, creating if necessary
func GetLoggerInstance() *log.Logger {
	once.Do(func() {
		logger = createLogger(filename)
	})
	return logger
}

func createLogger(fname string) *log.Logger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	return log.New(file, "amnisiac.", log.Lshortfile)

}
