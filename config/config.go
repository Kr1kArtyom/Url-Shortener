package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

type DatabaseConfig struct {
	DbURL string `yaml:"db_url"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

func Parse(configPath string) (*Config, error) {

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("can't get config path: %s", err.Error())
	}

	yamlConf, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("can't read conf: %s", err.Error())
	}

	var config Config

	err = yaml.Unmarshal(yamlConf, &config)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshall conf: %s", err.Error())
	}

	return &config, nil
}
