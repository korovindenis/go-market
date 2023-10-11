package config

import (
	"fmt"
	"os"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

const CONFIG_DEFAULT_PATH = "./configs/config.dev.yaml"

type config struct {
	webserver `koanf:"web_server"`
}

type webserver struct {
	mode     string `koanf:"web_server.mode"`
	host     string `koanf:"web_server.host"`
	port     int    `koanf:"web_server.port"`
	timeouts struct {
		idle       int `koanf:"web_server.timeouts.idle"`
		readHeader int `koanf:"web_server.timeouts.readHeader"`
		read       int `koanf:"web_server.timeouts.read"`
		write      int `koanf:"web_server.timeouts.write"`
	}
	maxHeaderBytes int `koanf:"web_server.maxHeaderBytes"`
}

func New() (*config, error) {
	k := koanf.New(".")
	configPath := CONFIG_DEFAULT_PATH
	if envVarValue := os.Getenv("CONFIG_PATH"); envVarValue != "" {
		configPath = envVarValue
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	var cfg config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	return &cfg, nil
}

func (c *config) GetServerMode() string {
	return c.webserver.mode
}

func (c *config) GetServerHost() string {
	return c.webserver.host
}

func (c *config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.webserver.host, c.webserver.port)
}

func (c *config) GetServerTimeoutIdle() time.Duration {
	return time.Duration(c.webserver.timeouts.idle) * time.Second
}

func (c *config) GetServerTimeoutReadHeader() time.Duration {
	return time.Duration(c.webserver.timeouts.readHeader) * time.Second
}

func (c *config) GetServerTimeoutRead() time.Duration {
	return time.Duration(c.webserver.timeouts.read) * time.Second
}

func (c *config) GetServerTimeoutWrite() time.Duration {
	return time.Duration(c.webserver.timeouts.write) * time.Second
}

func (c *config) GetServerMaxHeaderBytes() int {
	return c.webserver.maxHeaderBytes
}
