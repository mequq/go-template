package config

type (
	LogingConfig struct {
		Observability Observability `mapstructure:"observability"`
	}
	Observability struct {
		Logging Logging `mapstructure:"logging"`
		Otel    Otel    `mapstructure:"otel"`
	}
	Otel struct {
		Enabled bool   `mapstructure:"enabled"`
		Address string `mapstructure:"address"`
	}

	Logging struct {
		Level    string   `mapstructure:"level" `
		Logstash Logstash `mapstructure:"logstash"`
	}
	Logstash struct {
		Enabled bool   `mapstructure:"enabled"`
		Address string `mapstructure:"address"`
	}

	Server struct {
		HTTPServer HTTPServer `mapstructure:"http"`
	}

	HTTPServer struct {
		Port       int    `koanf:"Port"`
		Host       string `mapstructure:"host"`
		Production bool   `mapstructure:"production"`
	}
)
