package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/davidaparicio/ambari-to-opsgenie/api/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var v *viper.Viper

func init() {
	v = viper.New()
	v.AddConfigPath("/etc/appname/")  // path to look for the config file in
	v.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	v.AddConfigPath(".")              // optionally look for config in the working directory
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.WithError(err).Error("Fatal error config file")
		return
	}
	logLevel, err := log.ParseLevel(v.GetString("agent.loglevel"))
	if err == nil {
		log.SetLevel(logLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	var nbCritical, nbWarning int
	nbCritical, nbWarning = getNumbers()
	notifyBlinky(nbCritical, nbWarning)
	printBitbar(nbCritical, nbWarning)
}

func printBitbar(nbCritical int, nbWarning int) {
	//https://github.com/matryer/bitbar/issues/467#issuecomment-392062171
	//https://github.com/matryer/bitbar/issues/461#issuecomment-373689024
	message := "HADOOP: \033[31m" + strconv.Itoa(nbCritical) + "\033[34m / \033[33m" + strconv.Itoa(nbWarning) + "\033[0m | ansi=true"
	fmt.Println(message)
	fmt.Println("---")
	fmt.Println("Open a terminal | bash='ssh ${LOGIN}@${IP}'") //https://github.com/matryer/bitbar/pull/179
	fmt.Println("Refresh numbers | refresh=true")              //https://github.com/matryer/bitbar-plugins/blob/master/Tutorial/ansi.sh
	fmt.Println("---")
	fmt.Println("Web UIs")
	fmt.Println("Ambari | color=#123def href=https://localhost/gateway/default/ambari/")
	fmt.Println("YARN | color=#123def href=https://localhost/gateway/default/yarn/cluster/scheduler/")
	fmt.Println("Oozie | color=#123def href=https://localhost/gateway/default/oozie/")
	fmt.Println("Spark | color=#123def href=https://localhost/gateway/default/sparkhistory/")
	fmt.Println("LogSearch | color=#123def href=https://localhost/gateway/default/logsearch/")
	fmt.Println("---")
	fmt.Println("Dashboards")
	fmt.Println("Grafana | color=#123def href=https://grafana")
	fmt.Println("Kibana LDP | color=#123def href=https://kibana")
	fmt.Println("---")
	fmt.Println("Tooling")
	fmt.Println("-- Opsgenie | color=#123def href=https://myapp.app.eu.opsgenie.com/alert/list")
	//https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html
}

func getNumbers() (nbCritical, nbWarning int) {
	nbCritical = 0
	nbWarning = 0
	items, err := getAmbariAlert()
	if err != nil {
		log.WithError(err).Error("Fail to get Alert")
	}
	for _, item := range items {
		if item.Alert.State == "WARNING" {
			nbWarning++
		}
		if item.Alert.State == "CRITICAL" {
			nbCritical++
		}
	}
	return nbCritical, nbWarning
}

func getAmbariAlert() (alert []types.Item, err error) {
	u, err := url.Parse(v.GetString("ambari.url"))
	if err != nil {
		log.WithError(err).Error("parsing url")
		return
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ServerName: u.Host},
		},
	}
	req, err := http.NewRequest(http.MethodGet, v.GetString("ambari.url")+"/api/v1/clusters/prod/alerts?fields=*&Alert/maintenance_state=OFF", nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(v.GetString("ambari.username"), v.GetString("ambari.password"))
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

func notifyBlinky(nbCritical, nbWarning int) {
	if nbCritical != 0 {
		res, err := http.Get("https://app.getblinky.io/api/v1/ingest/webhook/5b0adf41-a91a-4e96-9265-f4081e0c30f4")
		if err != nil {
			log.WithError(err).Error("Fail to call Blinky critical webhook")
		} else {
			log.Info("Blinky webhook critical called\n")
			log.Infof("status code: %d\n", res.StatusCode)
		}
	} else if nbWarning != 0 {
		res, err := http.Get("https://app.getblinky.io/api/v1/ingest/webhook/c919e0fa-cbf2-4948-a111-a5dee3192d19")
		if err != nil {
			log.WithError(err).Error("Fail to call Blinky warning webhook")
		} else {
			log.Info("Blinky webhook warning called\n")
			log.Infof("status code: %d\n", res.StatusCode)
		}
	}
	/* FOR DEBUG else {
		res, err := http.Get("https://app.getblinky.io/api/v1/ingest/webhook/af2f8552-58a7-4a04-a43b-844887dcef8e")
		if err != nil {
			logrus.WithError(err).Error("Fail to call Blinky OK webhook")
		} else {
			logrus.Info("Blinky webhook OK called\n")
			logrus.Infof("status code: %d\n", res.StatusCode)
		}
	} */
}
