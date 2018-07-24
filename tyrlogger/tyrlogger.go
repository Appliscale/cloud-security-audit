package tyrlogger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	*log.Logger
	filename string
}

var logger *Logger
var once sync.Once

func GetInstance() *Logger {
	once.Do(func() {
		logger = createLogger("tyr.log")
	})
	return logger
}

func createLogger(fname string) *Logger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	return &Logger{
		filename: fname,
		Logger:   log.New(file, "Tyr ", log.Lshortfile),
	}
}
