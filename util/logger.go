package util

import (
	"github.com/sirupsen/logrus"
)

func (c Config) ConfigLogger() (err error) {
	// To get a timestamp format like DEBU[2022-01-11T23:15:56+01:00] Debug message
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05-07:00", FullTimestamp: true})

	logLevel, err := logrus.ParseLevel(c.V.GetString("agent.loglevel_unencrypted"))
	if err == nil {
		logrus.SetLevel(logLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return
}
