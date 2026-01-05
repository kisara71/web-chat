package logger

import (
	"log"
	"os"
	"sync"
)

var (
	once sync.Once
	lgr  *log.Logger
)

func L() *log.Logger {
	once.Do(func() {
		lgr = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	})
	return lgr
}
