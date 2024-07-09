package biz

import (
	"application/internal/biz/sample_entity"
	"github.com/google/wire"
)

var BizProviderSet = wire.NewSet(sample_entity.NewSampleEntity)
