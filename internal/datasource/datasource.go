package datasource

import (
	healthzrepo "application/internal/datasource/healthz"
	samplememrepo "application/internal/datasource/sample/memory"

	"github.com/google/wire"
)

var DataProviderSet = wire.NewSet(samplememrepo.NewSampleEntity, healthzrepo.NewHealthzDS)
