package lgr

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewDefaultLogger() *logrus.Logger {
	var log = logrus.New()
	setLogLevelByEnv(log)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})
	return log
}

func setLogLevelByEnv(log *logrus.Logger) {
	logLevel := strings.ToLower(os.Getenv("APP_LOG_LEVEL"))
	switch logLevel {
	case logrus.DebugLevel.String():
		log.SetLevel(logrus.DebugLevel)
	case logrus.InfoLevel.String():
		log.SetLevel(logrus.InfoLevel)
	case logrus.ErrorLevel.String():
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.ErrorLevel)
	}
}
