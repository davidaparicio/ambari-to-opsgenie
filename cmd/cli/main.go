package main

import (
	"context"
	"os"
	"time"

	"github.com/davidaparicio/ambari-to-opsgenie/internal"
	"github.com/davidaparicio/ambari-to-opsgenie/util"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/sirupsen/logrus"
)

const (
	EXIT_NOCONF_FILE = iota + 1
	EXIT_UNKNOWN_ERR
)

var l *logrus.Logger

func main() {
	var err error

	c := new(util.Config)
	err = util.LoadConfig(c)

	l = logrus.New()
	l.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05-07:00", FullTimestamp: true})
	l.SetLevel(logrus.DebugLevel)

	if err != nil {
		l.WithError(err).Error("cannot load config")
		os.Exit(EXIT_NOCONF_FILE)
	}

	l.Info(internal.CurrentVersion())
	c.AmbariOpgenieMapping = make(map[int]string)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go internal.SendHearbeat(ctx, c)

	c.AlertClient, err = alert.NewClient(&client.Config{
		ApiKey:         c.V.GetString("opsgenie.key"),
		OpsGenieAPIURL: client.API_URL_EU,
	})

	if err != nil {
		l.WithError(err).Error("Fail to Create Client")
		return
	}

	wait, _ := time.ParseDuration(c.V.GetString("ambari.interval_unencrypted"))
	ticker := time.NewTicker(wait)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			l.Debug("Stopping the Ambari-to-Opsgenie CLI..")
			return
		case <-ticker.C:
			items, err := internal.GetAmbariAlert(ctx, c)

			if err != nil {
				l.WithError(err).Error("Fail to get Alert")
				continue
			}

			for _, item := range items {
				//Check if the alert have been already created in Opsgenie
				opgenieID := c.AmbariOpgenieMapping[item.Alert.Id]

				if item.Alert.State == "OK" {
					if opgenieID == "" {
						//Nothing to do
						continue
					} else {
						//Closing the Opsgenie alert, because it's fixed
						if err = internal.CloseAlert(item.Alert, c); err != nil {
							l.WithError(err).Error("Fail to close Alert")
						}
					}
				}

				//item.Alert.State == "CRITICAL" or "WARNING"
				if opgenieID == "" {
					if err = internal.CreateAlert(item.Alert, c); err != nil {
						l.WithError(err).Error("Fail to send Alert")
					}
				} else {
					//Update the Opsgenie alert with a comment
					if err = internal.CommentAlert(item.Alert, c); err != nil {
						l.WithError(err).Error("Fail to comment Alert")
					}
				}
			}
		}
	}
}
