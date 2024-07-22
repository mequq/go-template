package datasource

import (
	memory2 "application/internal/v1/datasource/healthz/memory"
	"application/internal/v1/datasource/sample_entitiy/memory"
	"github.com/google/wire"
)

var DataProviderSet = wire.NewSet(memory.NewSampleEntity, memory2.NewHealthzDS)
