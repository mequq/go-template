package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config interface {
	Unmarshal(path string, cfg any) error
}

// config is the config.
type KoanfConfig struct {
	yamlConfigPath string

	k *koanf.Koanf
}

type ConfigOptions func(*KoanfConfig) error

func WithYamlConfigPath(path string) ConfigOptions {
	return func(c *KoanfConfig) error {
		c.yamlConfigPath = path
		return nil
	}
}

// newConfig creates a new config.
func NewKoanfConfig(opts ...ConfigOptions) (Config, error) {
	config := &KoanfConfig{
		k: koanf.New(""),
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
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "_", ".")
	}), nil); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *KoanfConfig) Unmarshal(path string, cfg any) error {
	return c.k.Unmarshal(path, cfg)
}
