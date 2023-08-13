package internal

import (
	"context"
	"time"

	"github.com/davidaparicio/ambari-to-opsgenie/util"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
	"github.com/sirupsen/logrus"
)

// SendHearbeat sends a heartbeat every x minutes.
// First you need to create a hearbeat in opsgenie and must be enabled
func SendHearbeat(ctx context.Context, c *util.Config) {

	HeartbeatClient, err := heartbeat.NewClient(&client.Config{
		ApiKey:         c.V.GetString("opsgenie.key"),
		OpsGenieAPIURL: client.API_URL_EU,
	})

	if err != nil {
		c.L.WithError(err).Error("Fail to Create heartbeat Client")
		return
	}

	wait, _ := time.ParseDuration(c.V.GetString("opsgenie.heartbeat.interval_unencrypted"))
	ticker := time.NewTicker(wait)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			c.L.Debug("sending Ping...")
			pingResult, err := HeartbeatClient.Ping(context.Background(), c.V.GetString("opsgenie.heartbeat.name_unencrypted"))
			if err != nil {
				c.L.WithError(err).Error("Fail to ping")
				continue
			}
			c.L.WithFields(logrus.Fields{
				"Message":      pingResult.Message,
				"ResponseTime": pingResult.ResponseTime,
			}).Debug("Pong received")
		}
	}
}
