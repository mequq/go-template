package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type ConfigInterface interface {
	Unmarshal(cfg any) error
}

// config is the config.
type Config struct {
	yamlConfigPath string

	k *koanf.Koanf
}

type ConfigOptions func(*Config) error

func WithYamlConfigPath(path string) ConfigOptions {
	return func(c *Config) error {
		c.yamlConfigPath = path
		return nil
	}
}

// newConfig creates a new config.
func NewKoanfConfig(opts ...ConfigOptions) (ConfigInterface, error) {
	config := &Config{
		yamlConfigPath: "config.yaml",
		k:              koanf.New(""),
	}

	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, err
		}
	}

	// load the config
	if config.yamlConfigPath != "" {
		if err := config.k.Load(file.Provider(config.yamlConfigPath), yaml.Parser()); err != nil {
			return nil, err
		}
	}

	//  load env variables
	if err := config.k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Unmarshal(cfg any) error {
	return c.k.Unmarshal("", cfg)
}
