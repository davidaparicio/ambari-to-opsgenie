package util

import (
	log "github.com/sirupsen/logrus"
)

// ConfigLogger configures the Logrus logger (timestamp+loglevel)
func (c Config) ConfigLogger() (err error) {
	// To get a timestamp format like DEBU[2022-01-11T23:15:56+01:00] Debug message
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02T15:04:05-07:00", FullTimestamp: true})

	logLevel, err := log.ParseLevel(c.V.GetString("agent.loglevel_unencrypted"))
	if err == nil {
		log.SetLevel(logLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
	return
}
