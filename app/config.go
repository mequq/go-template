package app

import (
	"errors"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type KConfig struct {
	*koanf.Koanf
}

var ErrFailedToLoadEnvVars = errors.New("failed to load env vars")

func NewKoanfConfig(runtimeFlags *runTimeFlags) (*KConfig, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(runtimeFlags.configYamlAddress), yaml.Parser()); err != nil {
		return nil, err
	}

	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "_", ".")
	}), nil); err != nil {
		return nil, errors.Join(ErrFailedToLoadEnvVars, err)
	}

	return &KConfig{k}, nil
}
