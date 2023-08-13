package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/davidaparicio/ambari-to-opsgenie/internal"
	"github.com/davidaparicio/ambari-to-opsgenie/util"
)

const (
	EXIT_NOCONF_FILE = iota + 1
)

var c *util.Config

func main() {
	var err error

	c = new(util.Config)
	err = util.LoadConfig(c)
	if err != nil {
		fmt.Println("cannot load config")
		os.Exit(EXIT_NOCONF_FILE)
	}

	c.L.Debugf("xbar/Ambari-to-Opsgenie %s", internal.Version)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var nbCritical, nbWarning int
	nbCritical, nbWarning = getNumbers(ctx, c)
	notifyBlinky(ctx, nbCritical, nbWarning)
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

func getNumbers(ctx context.Context, c *util.Config) (nbCritical, nbWarning int) {
	nbCritical = 0
	nbWarning = 0
	items, err := internal.GetAmbariAlert(ctx, c)
	if err != nil {
		c.L.WithError(err).Error("Fail to get Alert")
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

func notifyBlinky(ctx context.Context, nbCritical, nbWarning int) {
	if nbCritical != 0 {
		resp, err := httpGetWithContext(ctx, "https://app.getblinky.io/api/v1/ingest/webhook/5b0adf41-a91a-4e96-9265-f4081e0c30f4")
		if err != nil {
			c.L.WithError(err).Error("Fail to call Blinky critical webhook")
		} else {
			c.L.Info("Blinky webhook critical called\n")
			c.L.Infof("status code: %d\n", resp.StatusCode)
		}
	} else if nbWarning != 0 {
		resp, err := httpGetWithContext(ctx, "https://app.getblinky.io/api/v1/ingest/webhook/c919e0fa-cbf2-4948-a111-a5dee3192d19")
		if err != nil {
			c.L.WithError(err).Error("Fail to call Blinky warning webhook")
		} else {
			c.L.Info("Blinky webhook warning called\n")
			c.L.Infof("status code: %d\n", resp.StatusCode)
		}
	}
	/* FOR DEBUG else {
		res, err := httpGetWithContext(ctx, "https://app.getblinky.io/api/v1/ingest/webhook/af2f8552-58a7-4a04-a43b-844887dcef8e")
		if err != nil {
			c.L.WithError(err).Error("Fail to call Blinky OK webhook")
		} else {
			c.L.Info("Blinky webhook OK called\n")
			c.L.Infof("status code: %d\n", res.StatusCode)
		}
	} */
}

func httpGetWithContext(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)
	return resp, nil
}

// Deprecated: This function is normally not used by httpGetWithContext
// because http.NoBody.Close always returns nil
func closeResponse(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		c.L.WithError(err).Error("Fail to close http.NewRequest.body")
	}
}
