package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
)

type Logger struct {
	log logrus.FieldLogger
}

func New(level string, path string) (*Logger, error) {
	logger := logrus.New()

	loggerFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		err = fmt.Errorf("error during the opening log loggerFile: %w", err)
		return nil, err
	}

	logger.SetOutput(io.MultiWriter(os.Stdout, loggerFile))

	loggerLevel, err := logrus.ParseLevel(level)
	if err != nil {
		err = fmt.Errorf("error in parsing log level: %w", err)
		return nil, err
	}
	logger.SetLevel(loggerLevel)

	logger.SetFormatter(&prefixed.TextFormatter{DisableColors: true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})

	logger.WithFields(logrus.Fields{"level": level, "loggerFile": path}).Debug("Logger setup OK")

	return &Logger{logger}, nil
}

func (l Logger) Info(msg string) {
	l.log.Info(msg)
}

func (l Logger) Error(msg string) {
	l.log.Error(msg)
}

func (l Logger) Warn(msg string) {
	l.log.Warn(msg)
}

func (l Logger) Debug(msg string) {
	l.log.Debug(msg)
}
