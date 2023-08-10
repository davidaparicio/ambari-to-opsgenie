package internal

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/davidaparicio/ambari-to-opsgenie/api/types"
	"github.com/davidaparicio/ambari-to-opsgenie/util"
	log "github.com/sirupsen/logrus"
)

func GetAmbariAlert(c *util.Config) (alert []types.Item, err error) {

	u, err := url.Parse(c.V.GetString("ambari.url_unencrypted"))
	if err != nil {
		log.WithError(err).Error("parsing url")
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ServerName: u.Host},
		},
	}

	req, err := http.NewRequest(http.MethodGet, c.V.GetString("ambari.url_unencrypted"), nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(c.V.GetString("ambari.username_unencrypted"), c.V.GetString("ambari.password"))

	resp, err := client.Do(req)
	if err != nil {
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
		return
	}

	responseAlert := types.ResponseAlert{}

	err = json.Unmarshal(body, &responseAlert)
	if err != nil {
		return
	}

	alert = responseAlert.Items

	return

}
