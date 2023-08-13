package internal

import (
	"context"
	"strconv"

	"github.com/davidaparicio/ambari-to-opsgenie/api/types"
	"github.com/davidaparicio/ambari-to-opsgenie/util"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	log "github.com/sirupsen/logrus"
)

// CommentAlert comments an Opsgenie alert, using the Opsgenie Go SDK.
func CommentAlert(ambariAlert types.Alert, c *util.Config) (err error) {
	commentResult, err := c.AlertClient.AddNote(context.Background(), &alert.AddNoteRequest{
		IdentifierType:  alert.ALERTID,
		IdentifierValue: c.AmbariOpgenieMapping[ambariAlert.Id],
		Note:            ambariAlert.Text,
	})

	if err != nil {
		log.WithError(err).Error("Fail to comment Alert")
		return
	}

	status, err := commentResult.RetrieveStatus(context.Background())

	if !status.IsSuccess {
		log.Debug(status.Status)
		return
	}

	return
}

// CloseAlert closes an Opsgenie alert and remove it from the AmbariOpgenieMapping map.
func CloseAlert(ambariAlert types.Alert, c *util.Config) (err error) {
	closeResult, err := c.AlertClient.Close(context.Background(), &alert.CloseAlertRequest{
		IdentifierType:  alert.ALERTID,
		IdentifierValue: c.AmbariOpgenieMapping[ambariAlert.Id],
	})

	if err != nil {
		log.WithError(err).Error("Fail to Close Alert")
		return
	}

	status, err := closeResult.RetrieveStatus(context.Background())

	if !status.IsSuccess {
		log.Debug(status.Status)
		if status.Status != "Alert is already closed." {
			return
		}
	}

	delete(c.AmbariOpgenieMapping, ambariAlert.Id)

	return
}

// CreateAlert creates a new Opsgenie alert and save it into the AmbariOpgenieMapping map.
func CreateAlert(ambariAlert types.Alert, c *util.Config) (err error) {

	var priority alert.Priority
	switch ambariAlert.State {
	case "WARNING":
		priority = alert.P5
	case "CRITICAL":
		priority = alert.P3
	default:
		return
	}

	createResult, err := c.AlertClient.Create(context.Background(), &alert.CreateAlertRequest{
		Message:     ambariAlert.Label,
		Description: ambariAlert.Text,
		Entity:      ambariAlert.ServiceName,
		Source:      ambariAlert.DostName,
		Priority:    priority,
		Details: map[string]string{
			"id":            strconv.Itoa(ambariAlert.Id),
			"ComponentName": ambariAlert.ComponentName,
			"ClusterName":   ambariAlert.ClusterName,
			"Scope":         ambariAlert.Scope,
		},
	})

	if err != nil {
		log.WithError(err).Error("Fail to Create Alert")
		return
	}

	createStatus, err := createResult.RetrieveStatus(context.Background())

	if !createStatus.IsSuccess {
		log.Debug(createStatus.Status)
		return
	}

	c.AmbariOpgenieMapping[ambariAlert.Id] = createStatus.AlertID

	return

}
