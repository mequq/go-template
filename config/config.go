package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config is the config.
type Config struct {
	Tracing struct {
		Endpoint      string  `mapstructure:"endpoint"`
		ServiceName   string  `mapstructure:"servicename"`
		Environment   string  `mapstructure:"environment"`
		SamplingRatio float64 `mapstructure:"samplingratio"`
		TUrl          string  `mapstructure:"turl"`
	} `mapstructure:"tracing"`
	Server struct {
		// port is the port the server will listen on.
		Port string `mapstructure:"port"`
		// address is the address the server will listen on.
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`
	Redis struct {
		// host is the redis host.
		Host string `mapstructure:"host"`
		// port is the redis port.
		Port string `mapstructure:"port"`
		// password is the redis password.
		Password string `mapstructure:"password"`
		// db is the redis db.
		DB int `mapstructure:"db"`
	} `mapstructure:"redis"`
	Mongo struct {
		// host is the mongo host.
		Host string `mapstructure:"host"`
		// port is the mongo port.
		Port string `mapstructure:"port"`
		// username is the mongo username.
		Username string `mapstructure:"username"`
		// password is the mongo password.
		Password string `mapstructure:"password"`
		// db is the mongo db.
		DB string `mapstructure:"db"`
	} `mapstructure:"mongo"`
}

// NewConfig creates a new config.
func NewConfig() *Config {
	return &Config{}
}

// SetDefaults sets the default values.
func (c *Config) SetDefaults() {
	c.Server.Port = "8080"
	c.Server.Host = "0.0.0.0"
}

// Validate validates the config.
func (c *Config) Validate() error {
	return nil
}

// viper load config from path
func LoadConfig(configAddress string) (*Config, error) {
	// create a new config
	config := NewConfig()
	// create a new viper instance
	v := viper.New()
	// set the config name
	v.SetConfigFile(configAddress)
	// set env replacement
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//
	v.AutomaticEnv()
	// read the config
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	// unmarshal the config
	if err := v.Unmarshal(config); err != nil {
		return nil, err
	}
	// validate the config
	if err := config.Validate(); err != nil {
		return nil, err
	}
	// return the config
	return config, nil
}
