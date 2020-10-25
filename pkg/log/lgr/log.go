package lgr

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewDefaultLogger() *logrus.Logger {
	var log = logrus.New()
	setLogLevelByEnv(log)
	log.SetReportCaller(true)
	return log
}

func setLogLevelByEnv(log *logrus.Logger) {
	logLevel := os.Getenv("APP_LOG_LEVEL")
	switch logLevel {
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.ErrorLevel)
	}
}
