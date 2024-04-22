package logger

import (
	"fmt"
	"reflect"
	"time"

	"github.com/SashaMelva/car_catalog/internal/config"
	"go.uber.org/zap"
)

func New(conf *config.ConfigLogger, filePath string) *zap.SugaredLogger {
	fileName := time.Now()
	fmt.Println(conf.Level, reflect.TypeOf(conf.Level))

	logConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(conf.Level),
		DisableCaller:    true,
		Development:      true,
		Encoding:         conf.LogEncoding,
		OutputPaths:      []string{"stdout", filePath + fileName.Format("01-02-2006") + ".log"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	}

	logger := zap.Must(logConfig.Build()).Sugar()

	logger.Info("Started")
	logger.Debug("Debug mode enabled")
	return logger
}
