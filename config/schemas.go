package config

type (
	LogingConfig struct {
		Observability Observability `mapstructure:"observability"`
	}
	Observability struct {
		Logging Logging `mapstructure:"logging"`
	}
	Logging struct {
		Level    string   `mapstructure:"level" `
		Logstash Logstash `mapstructure:"logstash"`
	}
	Logstash struct {
		Enabled bool   `mapstructure:"enabled"`
		Address string `mapstructure:"address"`
	}
	HTTPServer struct {
		Port       int    `mapstructure:"port"`
		Host       string `mapstructure:"host"`
		Production bool   `mapstructure:"production"`
	}
)
