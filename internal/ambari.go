package internal

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/davidaparicio/ambari-to-opsgenie/api/types"
	"github.com/davidaparicio/ambari-to-opsgenie/util"
)

// GetAmbariAlert calls the Ambari API to retrieve all alerts of a Hadoop cluster
func GetAmbariAlert(ctx context.Context, c *util.Config) (alert []types.Item, err error) {

	u, err := url.Parse(c.V.GetString("ambari.url_unencrypted"))
	if err != nil {
		c.L.WithError(err).Error("parsing url")
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ServerName: u.Host, MinVersion: tls.VersionTLS12},
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.V.GetString("ambari.url_unencrypted"), nil)
	if err != nil {
		c.L.WithError(err).Error("creating the HTTP NewRequest")
		return
	}

	req.SetBasicAuth(c.V.GetString("ambari.username_unencrypted"), c.V.GetString("ambari.password"))

	resp, err := client.Do(req)
	if err != nil {
		c.L.WithError(err).Error("sending the HTTP request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 299 {
		err = errors.New(resp.Status)
		return
	}

	//read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.L.WithError(err).Error("reading the HTTP body")
		return
	}

	responseAlert := types.ResponseAlert{}

	err = json.Unmarshal(body, &responseAlert)
	if err != nil {
		c.L.WithError(err).Error("unmarshaling the JSON body")
		return
	}

	alert = responseAlert.Items
	return
}
