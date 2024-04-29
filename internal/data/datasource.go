package data

import (
	"application/config"
	"context"
	"errors"

	"log/slog"
)

// Erors
var (
	ErrorRepoNotInitialized = errors.New("repo not initialized")
)

// DataSouce is the struct that holds all the data sources
type DataSource struct {
	logger *slog.Logger
	cfg    config.ConfigInterface

	ctx context.Context
}

// NewDataSource creates a new DataSource
func NewDataSource(ctx context.Context, logger *slog.Logger, cfg config.ConfigInterface) (*DataSource, error) {
	ds := &DataSource{
		logger: logger.With("module", "repo"),
		cfg:    cfg,
		ctx:    ctx,
	}
	err := ds.Init()
	if err != nil {
		return nil, err
	}
	return ds, nil
}

type DataSourceConfig struct {
	Datasource struct {
		Mysql struct {
			Enabled               bool
			DSN                   string
			ConnectionPoolEnabled bool `koanf:"connection_pool_enabled"`
			ConnectionPoolMaxIdle int  `koanf:"connection_pool_max_idle"`
			ConnectionPoolMaxOpen int  `koanf:"connection_pool_max_open"`
		}
		Redis struct {
			Enabled  bool
			Address  string
			DB       int
			PassWord string
		}
	}
}

func (ds *DataSource) Init() error {

	cfg := &DataSourceConfig{}

	if err := ds.cfg.Unmarshal(cfg); err != nil {
		panic(err)
	}
	ds.logger.Info("config", "config", cfg)
	// err := ds.InitSQL()
	// if err != nil {
	// 	return err
	// }
	return nil
}

// func (ds *DataSource) Close() error {
// 	return nil
// }

// func (ds *DataSource) InitSQL() error {

// 	var err error
// 	dns := ds.cfg.DatasourceConfig.Mysql.Dns
// 	cfg := &gorm.Config{
// 		TranslateError: true,
// 	}
// 	ds.mysqlDB, err = gorm.Open(mysql.Open(dns), cfg)
// 	if err != nil {
// 		return err
// 	}
// 	sqlDB, err := ds.mysqlDB.DB()
// 	if err != nil {
// 		return err
// 	}

// 	sqlDB.SetMaxIdleConns(3)
// 	sqlDB.SetMaxOpenConns(100)
// 	sqlDB.SetConnMaxLifetime(time.Minute * 30)
// 	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
// 	return nil
// }
