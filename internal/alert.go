package internal

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/davidaparicio/ambari-to-opsgenie/api/types"
	"github.com/davidaparicio/ambari-to-opsgenie/util"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
)

var ErrOperationNotPermit = errors.New("can't create an Opsgenie alert for OK status")
var ErrOperationUnknown = errors.New("can't create an Opsgenie alert for an unknown status")

// CommentAlert comments an Opsgenie alert, using the Opsgenie Go SDK.
func CommentAlert(ctx context.Context, ambariAlert types.Alert, c *util.Config) (err error) {
	ctx, cancelled := context.WithCancel(ctx)
	defer cancelled()
	commentResult, err := c.AlertClient.AddNote(ctx, &alert.AddNoteRequest{
		IdentifierType:  alert.ALERTID,
		IdentifierValue: c.AmbariOpgenieMapping[ambariAlert.Id],
		Note:            ambariAlert.Text,
	})

	if err != nil {
		c.L.WithError(err).Error("Fail to comment Alert")
		return err
	}

	status, err := commentResult.RetrieveStatus(ctx)

	if !status.IsSuccess {
		c.L.Debug(status.Status)
		return err
	}

	return nil
}

// CloseAlert closes an Opsgenie alert and remove it from the AmbariOpgenieMapping map, using the Opsgenie Go SDK.
func CloseAlert(ctx context.Context, ambariAlert types.Alert, c *util.Config) (err error) {
	ctx, cancelled := context.WithCancel(ctx)
	defer cancelled()
	closeResult, err := c.AlertClient.Close(ctx, &alert.CloseAlertRequest{
		IdentifierType:  alert.ALERTID,
		IdentifierValue: c.AmbariOpgenieMapping[ambariAlert.Id],
	})

	if err != nil {
		c.L.WithError(err).Error("Fail to Close Alert")
		return err
	}

	status, err := closeResult.RetrieveStatus(ctx)

	if !status.IsSuccess {
		c.L.Debug(status.Status)
		if status.Status != "Alert is already closed." {
			c.L.WithError(err).Error("Fail to close a closed alert ;)")
			return err
		}
	}

	delete(c.AmbariOpgenieMapping, ambariAlert.Id)

	return nil
}

// CreateAlert creates a new Opsgenie alert and save it into the AmbariOpgenieMapping map, using the Opsgenie Go SDK.
func CreateAlert(ctx context.Context, ambariAlert types.Alert, c *util.Config) (err error) {

	var priority alert.Priority
	switch ambariAlert.State {
	case "WARNING":
		priority = alert.P5
	case "CRITICAL":
		priority = alert.P3
	case "OK":
		return OperationNotPossibleError("OK to close")
	default:
		return OperationUnknownError("CreateAlert(switch case = default)")
	}

	ctx, cancelled := context.WithCancel(ctx)
	defer cancelled()
	createResult, err := c.AlertClient.Create(ctx, &alert.CreateAlertRequest{
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
		c.L.WithError(err).Error("Fail to Create Alert")
		return err
	}

	createStatus, err := createResult.RetrieveStatus(ctx)

	if !createStatus.IsSuccess {
		c.L.Debug(createStatus.Status)
		return err
	}

	c.AmbariOpgenieMapping[ambariAlert.Id] = createStatus.AlertID

	return nil
}

func OperationNotPossibleError(op string) error {
	return fmt.Errorf("OperationNotPermit %w : %s", ErrOperationNotPermit, op)
}

func OperationUnknownError(op string) error {
	return fmt.Errorf("OperationUnknown %w : %s", ErrOperationUnknown, op)
}
