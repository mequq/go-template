package repo

import (
	"application/internal/biz"

	"github.com/google/wire"
)

var RepoProvider = wire.NewSet(
	NewHealthzDS,
	wire.Bind(new(biz.RepositoryHealthzer), new(*healthz)),
)
