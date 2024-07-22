package datasource

import (
	"github.com/google/wire"

	memory2 "application/internal/v1/datasource/healthz/memory"
	"application/internal/v1/datasource/sampleentity/memory"
)

var DataProviderSet = wire.NewSet(memory.NewSampleEntity, memory2.NewHealthzDS)
