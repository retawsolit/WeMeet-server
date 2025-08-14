package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Client struct {
		Path string `yaml:"path"`
		Port int    `yaml:"port"`
	} `yaml:"client"`
}

func Load(path string) (*AppConfig, error) {
	if path == "" {
		path = DefaultConfigFile
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrReadConfig
	}
	var c AppConfig
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, ErrParseConfig
	}

	// Áp fallback đơn giản nếu thiếu
	if c.Client.Port == 0 {
		c.Client.Port = DefaultClientPort
	}
	if c.Client.Path == "" {
		c.Client.Path = DefaultClientPath
	}
	return &c, nil
}
