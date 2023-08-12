package util

import (
	"bytes"
	"os"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mozilla.org/sops/v3/decrypt"
)

// A Config specifies the configuration for alerts.
type Config struct {
	V                    *viper.Viper
	L                    *logrus.Logger
	AmbariOpgenieMapping map[int]string
	AlertClient          *alert.Client
}

const ageKeyFile = "secrets/age.key"

// LoadConfig reads encrypted configuration from file or environment variables.
func LoadConfig(c *Config) (err error) {
	//https://blog.gitguardian.com/a-comprehensive-guide-to-sops/
	if _, ok := os.LookupEnv("SOPS_AGE_KEY_FILE"); !ok {
		// Set the hardcoded constant of ageKeyFile
		err = os.Setenv("SOPS_AGE_KEY_FILE", ageKeyFile)
		if err != nil {
			c.L.WithError(err).Error("Fail to set the OS environment variable")
		}
	}

	//Decrypt, like https://github.com/dailymotion-oss/octopilot/blob/280196f325b8051315e40170ab786355ea856e14/update/sops/sops_test.go
	actualCleartextData, err := decrypt.File("configs/config.enc.yaml", "yaml")
	if err != nil {
		c.L.WithError(err).Error("Decrypt config file fail")
	}

	c.V = viper.New()
	c.V.SetConfigType("yaml")
	if err = c.V.ReadConfig(bytes.NewBuffer(actualCleartextData)); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			c.L.WithError(err).Error("Config file not found. Fail to init config")
		} else {
			c.L.WithError(err).Error("Config file found but viper fails to ReadConfig")
		}
	}
	/*if err = c.ConfigLogger(); err != nil {
		fmt.Printf("[ERRO] to configure the logrus logger...\n%v\n", err)
	}*/
	return
}
