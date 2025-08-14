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

	DatabaseInfo struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Charset  string `yaml:"charset"`
	} `yaml:"database_info"`

	RedisInfo struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		DB   int    `yaml:"db"`
	} `yaml:"redis_info"`

	NatsInfo struct {
		NatsURLs   []string `yaml:"nats_urls"`
		NatsWSURLs []string `yaml:"nats_ws_urls"`
	} `yaml:"nats_info"`

	LivekitInfo struct {
		Host   string `yaml:"host"`
		APIKey string `yaml:"api_key"`
		Secret string `yaml:"secret"`
	} `yaml:"livekit_info"`

	Etherpad struct {
		Host string `yaml:"host"`
	} `yaml:"etherpad"`
}

func Load(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c AppConfig
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
