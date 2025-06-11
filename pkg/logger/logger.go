package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"

	Library "mini-accounting/library"
)

var logger *logrus.Logger

func New(
	library Library.Library,
) {
	logLevel := logrus.DebugLevel
	log := logrus.New()
	log.SetLevel(logLevel)
	rotateFileHook, err := library.GetNewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename: fmt.Sprintf("logs/%s.log", library.GetNow().Format("2006-01-02")),
		MaxSize:  50, // MB
		MaxAge:   28, // DAYS
		Level:    logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat:   "2006-01-02 15:04:05",
			PrettyPrint:       true,
			DisableHTMLEscape: true,
		},
	})

	if err != nil {
		logrus.Panic(err)
	}

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		PrettyPrint:       true,
		DisableHTMLEscape: true,
	})

	log.AddHook(rotateFileHook)

	logger = log
}

func WriteLog(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

func GetLogger() *logrus.Logger {
	return logger
}
