package config

import (
	"fmt"
	"os"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/spf13/pflag"
)

const configDefaultPath = "./Configs/Config.dev.yaml"

type config struct {
	App        `koanf:"app"`
	Httpserver `koanf:"http_server"`
	Storage    `koanf:"storage"`
	Accrual    `koanf:"accrual"`
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

type Accrual struct {
	Address string `koanf:"address"`
}

func New() (*config, error) {
	k := koanf.New(".")
	configPath := configDefaultPath
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
	address := pflag.StringP("address", "a", "", "Address and port to run the service")
	pflag.Parse()

	addr := *address

	if addr == "" {
		addr = os.Getenv("RUN_ADDRESS")
	}

	if c.Httpserver.Host != "" && c.Httpserver.Port != 0 {
		addr = fmt.Sprintf("%s:%d", c.Httpserver.Host, c.Httpserver.Port)
	}

	if addr == "" {
		addr = "localhost:8082"
	}

	return addr
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
	address := pflag.StringP("database", "d", "", "Database connection string")
	pflag.Parse()

	addr := *address

	if addr == "" {
		addr = os.Getenv("DATABASE_URI")
	}

	if c.Storage.ConnectionString != "" {
		addr = c.Storage.ConnectionString
	}

	if addr == "" {
		addr = "host=127.0.0.1 user=go password=go dbname=go sslmode=disable"
	}

	return addr
}

func (c *config) GetStorageSalt() string {
	return c.Storage.Salt
}

func (c *config) GetAccrualAddress() string {
	address := pflag.StringP("accural", "r", "", "Accural service address")
	pflag.Parse()

	addr := *address

	if addr == "" {
		addr = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	}

	if c.Accrual.Address != "" {
		addr = c.Accrual.Address
	}

	if addr == "" {
		addr = "http://localhost:8080"
	}

	return addr
}
