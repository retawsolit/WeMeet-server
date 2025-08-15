package config

import (
	"os"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type AppConfig struct {
	Client   ClientConfig `yaml:"client"`
	Database DatabaseInfo `yaml:"database_info"`
	Redis    RedisInfo    `yaml:"redis_info"`
	NATS     NATSInfo     `yaml:"nats_info"`
	LiveKit  LiveKitInfo  `yaml:"livekit_info"`
	Etherpad EtherpadInfo `yaml:"etherpad"`

	// Runtime connections - không serialize
	DB       *gorm.DB      `yaml:"-"`
	RDS      *redis.Client `yaml:"-"`
	NatsConn *nats.Conn    `yaml:"-"`
}

type ClientConfig struct {
	Path string `yaml:"path"`
	Port int    `yaml:"port"`
}

type DatabaseInfo struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Charset  string `yaml:"charset"`
}

type RedisInfo struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   int    `yaml:"db"`
}

type NATSInfo struct {
	NatsUrls   []string `yaml:"nats_urls"`
	NatsWsUrls []string `yaml:"nats_ws_urls"`
}

type LiveKitInfo struct {
	Host   string `yaml:"host"`
	ApiKey string `yaml:"api_key"`
	Secret string `yaml:"secret"`
}

type EtherpadInfo struct {
	Host string `yaml:"host"`
}

// Singleton pattern
var globalConfig *AppConfig

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

	// Set defaults
	if c.Client.Port == 0 {
		c.Client.Port = DefaultClientPort
	}
	if c.Client.Path == "" {
		c.Client.Path = DefaultClientPath
	}

	globalConfig = &c
	return &c, nil
}

func GetConfig() *AppConfig {
	return globalConfig
}
