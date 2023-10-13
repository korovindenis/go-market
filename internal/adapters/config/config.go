package config

import (
	"fmt"
	"os"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

const config_default_Path = "./Configs/Config.dev.yaml"

type Config struct {
	LogsLevel  string `koanf:"logs_level"`
	Httpserver `koanf:"http_server"`
	Storage    `koanf:"storage"`
}

type Httpserver struct {
	Mode           string `koanf:"mode"`
	Host           string `koanf:"host"`
	Port           int    `koanf:"port"`
	MaxHeaderBytes int    `koanf:"maxHeaderBytes"`
	Timeouts       struct {
		Idle       int `koanf:"idle"`
		ReadHeader int `koanf:"readHeader"`
		Read       int `koanf:"read"`
		Write      int `koanf:"write"`
	}
}

type Storage struct {
	ConnectionString string `koanf:"connection_string"`
}

func New() (*Config, error) {
	k := koanf.New(".")
	configPath := config_default_Path
	if envVarValue := os.Getenv("CONFIG_PATH"); envVarValue != "" {
		configPath = envVarValue
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("сonfig file does not exist: %s", configPath)
	}

	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading сonfig: %v", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error parsing сonfig: %v", err)
	}

	return &cfg, nil
}

func (c *Config) GetLogsLevel() string {
	return c.LogsLevel
}

func (c *Config) GetServerMode() string {
	return c.Httpserver.Mode
}

func (c *Config) GetServerHost() string {
	return c.Httpserver.Host
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Httpserver.Host, c.Httpserver.Port)
}

func (c *Config) GetServerTimeoutIdle() time.Duration {
	return time.Duration(c.Httpserver.Timeouts.Idle) * time.Second
}

func (c *Config) GetServerTimeoutReadHeader() time.Duration {
	return time.Duration(c.Httpserver.Timeouts.ReadHeader) * time.Second
}

func (c *Config) GetServerTimeoutRead() time.Duration {
	return time.Duration(c.Httpserver.Timeouts.Read) * time.Second
}

func (c *Config) GetServerTimeoutWrite() time.Duration {
	return time.Duration(c.Httpserver.Timeouts.Write) * time.Second
}

func (c *Config) GetServerMaxHeaderBytes() int {
	return c.Httpserver.MaxHeaderBytes
}

func (c *Config) GetStorageConnectionString() string {
	return c.Storage.ConnectionString
}
