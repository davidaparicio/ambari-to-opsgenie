package main

import (
	"os"
	"time"

	"github.com/davidaparicio/ambari-to-opsgenie/internal"
	"github.com/davidaparicio/ambari-to-opsgenie/util"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"

	log "github.com/sirupsen/logrus"
)

const (
	EXIT_NOCONF_FILE = iota + 1
	EXIT_UNKNOWN_ERR
)

func main() {
	var err error

	c := new(util.Config)
	err = util.LoadConfig(c)

	if err != nil {
		log.WithError(err).Error("cannot load config")
		os.Exit(EXIT_NOCONF_FILE)
	}

	go internal.SendHearbeat(c)

	c.AmbariOpgenieMapping = make(map[int]string)
	c.AlertClient, err = alert.NewClient(&client.Config{
		ApiKey:         c.V.GetString("opsgenie.key"),
		OpsGenieAPIURL: client.API_URL_EU,
	})

	if err != nil {
		log.WithError(err).Error("Fail to Create Client")
		return
	}

	wait, _ := time.ParseDuration(c.V.GetString("ambari.interval_unencrypted"))
	ticker := time.NewTicker(wait)

	for range ticker.C {
		items, err := internal.GetAmbariAlert(c)
		if err != nil {
			log.WithError(err).Error("Fail to get Alert")
			continue
		}

		for _, item := range items {
			opgenieID := c.AmbariOpgenieMapping[item.Alert.Id]

			if item.Alert.State == "OK" {
				if opgenieID == "" {
					continue
				} else {
					if err = internal.CloseAlert(item.Alert, c); err != nil {
						log.WithError(err).Error("Fail to close Alert")
					}
				}
			}

			if opgenieID == "" {
				if err = internal.CreateAlert(item.Alert, c); err != nil {
					log.WithError(err).Error("Fail to send Alert")
				}
			} else {

				//if err = checkAndUpdateAlert(item.Alert); err != nil {
				if err = internal.CommentAlert(item.Alert, c); err != nil {
					log.WithError(err).Error("Fail to update Alert")
				}
			}
		}
	}
}
