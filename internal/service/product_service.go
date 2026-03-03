package service

import (
	"context"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var CTX *context.Context
var Logger zap.Logger
var MTX *sync.Mutex
var RMTX *sync.RWMutex

func InitLogger() (*os.File, *zap.Logger) {
	time := time.Now().Format("2006-01-02 -- 15-04-05")

	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile("./logs/"+time+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}

	fileSyncer := zapcore.AddSync(logFile)

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewCore(
		encoder,
		fileSyncer,
		zap.InfoLevel,
	)

	logger := zap.New(core)
	Logger = *logger

	Logger.Info("Logger initialized successfully!")
	return logFile, logger
}
