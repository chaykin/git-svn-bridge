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

	AuthorsFile    string `yaml:"authorsFile"`
	LogWindowsSize uint   `yaml:"logWindowsSize"`

	SystemGitUserName   string `yaml:"systemGitUserName"`
	OsGitUserName       string `yaml:"osGitUserName"`
	CommitMessageFormat string `yaml:"commitMessageFormat"`

	LogFile string `yaml:"logFile"`
}

func GetConfig() *Config {
	if config == nil {
		confFile, err := os.Open("config.yml")
		if err != nil {
			panic(fmt.Errorf("could not read config file: %w", err))
		}

		defer closeFile(confFile)

		if err := yaml.NewDecoder(confFile).Decode(&config); err != nil {
			panic(fmt.Errorf("could not parse config file: %w", err))
		}
	}

	return config
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		panic(fmt.Errorf("could not close config file: %w", err))
	}
}
