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
// SOPS guide: https://blog.gitguardian.com/a-comprehensive-guide-to-sops/
func LoadConfig(c *Config) (err error) {
	// Set logrus logger, before using it
	c.ConfigLogger()

	// SOPS/AGE preparation, check SOPS_AGE_KEY_FILE OS variable
	if _, ok := os.LookupEnv("SOPS_AGE_KEY_FILE"); !ok {
		// Set the hardcoded constant of ageKeyFile
		err = os.Setenv("SOPS_AGE_KEY_FILE", ageKeyFile)
		if err != nil {
			c.L.WithError(err).Error("Fail to set the OS environment variable")
		}
	}

	// Decrypt using AGE key
	actualCleartextData, err := decrypt.File("configs/config.enc.yaml", "yaml")
	if err != nil {
		c.L.WithError(err).Error("Decrypt config file fail")
	}

	// Read decrypted configuration with Viper
	c.V = viper.New()
	c.V.SetConfigType("yaml")
	if err = c.V.ReadConfig(bytes.NewBuffer(actualCleartextData)); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			c.L.WithError(err).Error("Config file not found. Fail to init config")
		} else {
			c.L.WithError(err).Error("Config file found but viper fails to ReadConfig")
		}
	}

	// Configure the correct logrus level
	logLevel, err := logrus.ParseLevel(c.V.GetString("loglevel_unencrypted"))
	if err == nil {
		c.L.SetLevel(logLevel)
	} else {
		c.L.SetLevel(logrus.DebugLevel)
	}
	return
}
