package config

import (
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

func LoadConfig[C any](configFile string) (*C, error) {
	var cfg C
	if err := loadConfigFile(configFile, &cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadConfigFile[C any](configFile string, cfg *C) error {
	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	ymlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(ymlFile, &cfg)
}
