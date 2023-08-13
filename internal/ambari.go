package internal

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/davidaparicio/ambari-to-opsgenie/api/types"
	"github.com/davidaparicio/ambari-to-opsgenie/util"
)

const ERROR_STATUS_CODE = 299

var ErrStatusCode = errors.New("HTTP Error")

// GetAmbariAlert calls the Ambari API to retrieve all alerts of a Hadoop cluster
func GetAmbariAlert(ctx context.Context, c *util.Config) (alert []types.Item, err error) {
	url, err := url.Parse(c.V.GetString("ambari.url_unencrypted"))
	if err != nil {
		c.L.WithError(err).Error("parsing url")
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ServerName: url.Host, MinVersion: tls.VersionTLS12},
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.V.GetString("ambari.url_unencrypted"), nil)
	if err != nil {
		c.L.WithError(err).Error("creating the HTTP NewRequest")
		return nil, err
	}

	req.SetBasicAuth(c.V.GetString("ambari.username_unencrypted"), c.V.GetString("ambari.password"))

	resp, err := client.Do(req)
	if err != nil {
		c.L.WithError(err).Error("sending the HTTP request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= ERROR_STATUS_CODE {
		return nil, httpError(resp.Status)
	}

	//read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.L.WithError(err).Error("reading the HTTP body")
		return nil, err
	}

	responseAlert := types.ResponseAlert{}

	err = json.Unmarshal(body, &responseAlert)
	if err != nil {
		c.L.WithError(err).Error("unmarshaling the JSON body")
		return nil, err
	}

	alert = responseAlert.Items
	return
}

func httpError(statusCode string) error {
	return fmt.Errorf("OperationUnknown %w : %s", ErrStatusCode, statusCode)
}
