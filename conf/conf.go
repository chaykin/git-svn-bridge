package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var config *Config

type Config struct {
	ReposRoot   string `yaml:"reposRoot"`
	DbRoot      string `yaml:"dbRoot"`
	DbCacheSize uint64 `yaml:"dbCacheSize"`
}

func GetConfig() *Config {
	if config == nil {
		confFile, err := os.Open("config.yml")
		if err != nil {
			panic(fmt.Errorf("could not read config file: %w", err))
		}

		err = yaml.NewDecoder(confFile).Decode(&config)
		if err != nil {
			panic(fmt.Errorf("could not parse config file: %w", err))
		}

		err = confFile.Close()
		if err != nil {
			panic(fmt.Errorf("could not close config file: %w", err))
		}
	}

	return config
}
