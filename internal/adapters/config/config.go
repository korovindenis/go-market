package config

import (
	"fmt"
	"os"
	"time"

	"flag"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

const configDefaultPath = "./configs/config.dev.yaml"

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
	Address        string `koanf:"address"`
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
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	var config config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	config.parseFlags()

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

func (c *config) parseFlags() {
	flag.StringVar(&c.Httpserver.Address, "a", "localhost:8080", "Address and port to run the service")
	flag.StringVar(&c.Storage.ConnectionString, "d", "host=127.0.0.1 user=go password=go dbname=go sslmode=disable", "Database connection string")
	flag.StringVar(&c.Accrual.Address, "r", "http://localhost:8082", "Accural service address")
	flag.Parse()

	if envKey, err := getEnvVariable("RUN_ADDRESS"); err == nil {
		c.Httpserver.Address = envKey
	}
	if envKey, err := getEnvVariable("DATABASE_URI"); err == nil {
		c.Storage.ConnectionString = envKey
	}
	if envKey, err := getEnvVariable("ACCRUAL_SYSTEM_ADDRESS"); err == nil {
		c.Accrual.Address = envKey
	}
}

func (c *config) GetServerAddress() string {
	return c.Httpserver.Address
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

func (c *config) GetAccrualAddress() string {
	return c.Accrual.Address
}

func getEnvVariable(varName string) (string, error) {
	if envVarValue, exists := os.LookupEnv(varName); exists && envVarValue != "" {
		return envVarValue, nil
	}
	return "", entity.ErrEnvVarNotFound
}
