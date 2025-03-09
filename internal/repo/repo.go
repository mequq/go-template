package repo

import (
	healthzrepo "application/internal/repo/healthz"
	samplerepo "application/internal/repo/sample/memory"

	"github.com/google/wire"
)

var RepoProvider = wire.NewSet(healthzrepo.NewHealthzDS, samplerepo.NewSampleEntity)
