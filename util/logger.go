package util

import "github.com/sirupsen/logrus"

// ConfigLogger configures the Logrus logger (timestamp+loglevel).
func (c *Config) ConfigLogger() {
	c.L = logrus.New()
	// To get a timestamp format like DEBU[2022-01-11T23:15:56+01:00] Debug message
	c.L.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05-07:00",
		FullTimestamp:   true,
		// ForceColors:     true,
		// DisableColors:   false,
	})
	// Set logrus show line number
	// c.L.SetReportCaller(true)
}
