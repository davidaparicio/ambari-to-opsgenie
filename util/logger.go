package util

/*type myFormatter struct {
	logrus.TextFormatter
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 31 //gray
	case logrus.WarnLevel:
		levelColor = 33 //yellow
	case logrus.ErrorLevel:
		levelColor = 30 //gray
	default:
		levelColor = 36 //blue
	}
	return []byte(fmt.Sprintf("[%s] - \x1b[%dm%s\x1b[0m - %s\n", entry.Time.Format(f.TimestampFormat), levelColor, strings.ToUpper(entry.Level.String()), entry.Message)), nil
}*/

// ConfigLogger configures the Logrus logger (timestamp+loglevel)
/*func (c Config) ConfigLogger() (err error) {
	c.L = logrus.New()
	// To get a timestamp format like DEBU[2022-01-11T23:15:56+01:00] Debug message
	//c.L.SetFormatter(&myFormatter{logrus.TextFormatter{
	c.L.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05-07:00",
		FullTimestamp:   true,
		ForceColors:     true,
		DisableColors:   false,
	})

	logLevel, err := logrus.ParseLevel(c.V.GetString("loglevel_unencrypted"))
	if err == nil {
		c.L.SetLevel(logLevel)
	} else {
		c.L.SetLevel(logrus.DebugLevel)
	}
	//c.L.SetLevel(logrus.DebugLevel)

	//Set logrus show line number
	//c.L.SetReportCaller(true)
	return
}*/
