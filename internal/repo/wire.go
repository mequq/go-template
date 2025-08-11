package repo

import (
	"application/internal/biz"

	"github.com/google/wire"
)

var RepoProvider = wire.NewSet(
	NewHealthzDS,
	NewSampleEntity,
	wire.Bind(new(biz.RepositoryHealthzer), new(*healthz)),
	wire.Bind(new(biz.RepositorySampleer), new(*sampleEntity)),
)
