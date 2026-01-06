package logger

import (
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

var (
	once   sync.Once
	mu     sync.RWMutex
	global Logger
	writer io.Writer
)

func L() Logger {
	initDefault()
	mu.RLock()
	defer mu.RUnlock()
	return global
}

func Writer() io.Writer {
	initDefault()
	mu.RLock()
	defer mu.RUnlock()
	return writer
}

func Set(l Logger, w io.Writer) {
	if l == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	global = l
	if w != nil {
		writer = w
	}
}

func initDefault() {
	once.Do(func() {
		zapLogger, _ := zap.NewProduction()
		global = zapLogger.Sugar()
		writer = zapcore.AddSync(os.Stdout)
	})
}
