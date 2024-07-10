package datasource

import (
	memory2 "application/internal/datasource/healthz/memory"
	"application/internal/datasource/sample_entitiy/memory"
	"github.com/google/wire"
)

var DataProviderSet = wire.NewSet(memory.NewSampleEntity, memory2.NewHealthzDS)
