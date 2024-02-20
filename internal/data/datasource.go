package data

import (
	"context"
	"errors"

	"application/config"

	"log/slog"
)

// Erors
var (
	ErrorRepoNotInitialized = errors.New("repo not initialized")
)

// DataSouce is the struct that holds all the data sources
type DataSource struct {
	logger *slog.Logger
	cfg    *config.ViperConfig

	ctx context.Context
}

// NewDataSource creates a new DataSource
func NewDataSource(ctx context.Context, logger *slog.Logger, cfg *config.ViperConfig) (*DataSource, error) {
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

func (ds *DataSource) Init() error {
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
