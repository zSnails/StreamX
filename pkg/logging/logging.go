package logging

import (
	"github.com/sirupsen/logrus"
)

var logger logrus.Entry

func init() {
	logger = *logrus.NewEntry(logrus.StandardLogger())
    logger.Logger.SetLevel(logrus.DebugLevel)
	logger.Logger.SetFormatter(&logrus.TextFormatter{})
}

func Get() *logrus.Entry {
	return &logger
}
