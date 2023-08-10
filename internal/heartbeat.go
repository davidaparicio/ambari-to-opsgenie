package internal

import (
	"context"
	"time"

	"github.com/davidaparicio/ambari-to-opsgenie/util"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
	log "github.com/sirupsen/logrus"
)

// sendHearbeat sends a heartbeat every x minutes
// first you need to create a hearbeat in opsgenie and must be enabled
func SendHearbeat(c *util.Config) {

	HeartbeatClient, err := heartbeat.NewClient(&client.Config{
		ApiKey:         c.V.GetString("opsgenie.key"),
		OpsGenieAPIURL: client.API_URL_EU,
	})

	if err != nil {
		log.WithError(err).Error("Fail to Create heartbeat Client")
		return
	}

	wait, _ := time.ParseDuration(c.V.GetString("opsgenie.heartbeat.interval_unencrypted"))
	ticker := time.NewTicker(wait)

	for range ticker.C {
		log.Debug("sending Ping...")
		pingResult, err := HeartbeatClient.Ping(context.Background(), c.V.GetString("opsgenie.heartbeat.name_unencrypted"))
		if err != nil {
			log.WithError(err).Error("Fail to ping")
			continue
		}
		log.WithFields(log.Fields{
			"Message":      pingResult.Message,
			"ResponseTime": pingResult.ResponseTime,
		}).Debug("Pong received")

	}

}
