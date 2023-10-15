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

type config struct {
	App        `koanf:"app"`
	Httpserver `koanf:"http_server"`
	Storage    `koanf:"storage"`
}

type App struct {
	LogsLevel     string `koanf:"logs_level"`
	SecretKey     string `koanf:"secret_key"`
	TokenName     string `koanf:"token_name"`
	TokenLifeTime int    `koanf:"token_lifetime"`
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
	Salt             string `koanf:"salt"`
}

func New() (*config, error) {
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

	var config config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("error parsing сonfig: %v", err)
	}

	return &config, nil
}

func (c *config) GetTokenName() string {
	return c.App.TokenName
}

func (c *config) GetLogsLevel() string {
	return c.App.LogsLevel
}

func (c *config) GetAppSecretKey() string {
	return c.App.SecretKey
}

func (c *config) GetTokenLifeTime() time.Duration {
	return time.Duration(c.App.TokenLifeTime)
}

func (c *config) GetServerMode() string {
	return c.Httpserver.Mode
}

func (c *config) GetServerHost() string {
	return c.Httpserver.Host
}

func (c *config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Httpserver.Host, c.Httpserver.Port)
}

func (c *config) GetServerTimeoutIdle() time.Duration {
	return time.Duration(c.Httpserver.Timeouts.Idle) * time.Second
}

func (c *config) GetServerTimeoutReadHeader() time.Duration {
	return time.Duration(c.Httpserver.Timeouts.ReadHeader) * time.Second
}

func (c *config) GetServerTimeoutRead() time.Duration {
	return time.Duration(c.Httpserver.Timeouts.Read) * time.Second
}

func (c *config) GetServerTimeoutWrite() time.Duration {
	return time.Duration(c.Httpserver.Timeouts.Write) * time.Second
}

func (c *config) GetServerMaxHeaderBytes() int {
	return c.Httpserver.MaxHeaderBytes
}

func (c *config) GetStorageConnectionString() string {
	return c.Storage.ConnectionString
}

func (c *config) GetStorageSalt() string {
	return c.Storage.Salt
}
