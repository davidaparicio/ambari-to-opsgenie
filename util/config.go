package util

import (
	"bytes"
	"os"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mozilla.org/sops/v3/decrypt"
)

// A Config specifies the configuration for alerts.
// [TODO] The zero value for Config is a ready-to-use default configuration.
type Config struct {
	V                    *viper.Viper
	AmbariOpgenieMapping map[int]string
	AlertClient          *alert.Client
}

/*var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)*/

// https://github.com/golang/go/blob/457721cd52008146561c80d686ce1bb18285fe99/src/go/types/api.go#L110

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(c *Config) (err error) {
	//https://blog.gitguardian.com/a-comprehensive-guide-to-sops/
	const (
		ageKeyFile   = "secrets/age.key"
		agePublicKey = "age10qrzk97esg3lr6pcrgpk82p3sx3k2dqg477ufgx2l3gg564d7gjq2ux9vy"
	)
	err = os.Setenv("SOPS_AGE_KEY_FILE", ageKeyFile)
	if err != nil {
		log.WithError(err).Error("Fail to set the OS environment variable")
	}

	//Decrypt, like https://github.com/dailymotion-oss/octopilot/blob/280196f325b8051315e40170ab786355ea856e14/update/sops/sops_test.go
	actualCleartextData, err := decrypt.File("configs/config.enc.yaml", "yaml")
	if err != nil {
		log.WithError(err).Error("Decrypt config file fail")
	}

	c.V = viper.New()
	c.V.SetConfigType("yaml")
	// v.AddConfigPath(".") // err := v.ReadInConfig()
	if err = c.V.ReadConfig(bytes.NewBuffer(actualCleartextData)); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.WithError(err).Error("Config file not found. Fail to init config")
		} else {
			log.WithError(err).Error("Config file found but viper fails to ReadConfig")
		}
	}
	return
}
