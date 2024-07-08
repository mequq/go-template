package sample_entitiy

import "application/internal/entity"

type DataSource interface {
	Create() (id uint64, err error)
	Update() error
	List() []*entity.SampleEntity
	Delete() error
}
