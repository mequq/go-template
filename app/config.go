package app

import (
	"errors"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type kConfig struct {
	koanf.Koanf
}

func NewKoanfConfig(runtimeFlags *runTimeFlags) (*kConfig, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(runtimeFlags.configYamlAddress), yaml.Parser()); err != nil {
		return nil, err
	}

	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "_", ".")
	}), nil); err != nil {
		return nil, errors.Join(errors.New("failed to load env vars"), err)
	}

	return &kConfig{*k}, nil
}
