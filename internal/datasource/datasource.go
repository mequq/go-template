package datasource

import (
	"application/internal/datasource/sample_entitiy/memory"
	"github.com/google/wire"
)

var DataProviderSet = wire.NewSet(memory.NewSampleEntity)
