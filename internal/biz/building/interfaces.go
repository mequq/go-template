package buildingusecase

import (
	"context"
	"log/slog"
)

type BuildingUseCaseInterface interface {
	CreateBuilding(ctx context.Context, building *Building) error
	GetBuilding(ctx context.Context, id int) (*Building, error)
	GetBuildings(ctx context.Context) ([]*Building, error)
	UpdateBuilding(ctx context.Context, building *Building) error
	DeleteBuilding(ctx context.Context, id int) error
}

type BuildingUcas struct {
	logger *slog.Logger
}

func NewBuildingUseCase(logger *slog.Logger) BuildingUseCaseInterface {
	return &BuildingUcas{
		logger: logger,
	}
}

// create building
func (b *BuildingUcas) CreateBuilding(ctx context.Context, building *Building) error {
	return nil
}

// get building
func (b *BuildingUcas) GetBuilding(ctx context.Context, id int) (*Building, error) {
	return nil, nil
}

// get buildings
func (b *BuildingUcas) GetBuildings(ctx context.Context) ([]*Building, error) {
	return nil, nil
}

// update building
func (b *BuildingUcas) UpdateBuilding(ctx context.Context, building *Building) error {
	return nil
}

// delete building
func (b *BuildingUcas) DeleteBuilding(ctx context.Context, id int) error {
	return nil
}
