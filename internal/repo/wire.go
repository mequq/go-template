package repo

import (
	"application/internal/biz"

	"github.com/google/wire"
)

var RepoProvider = wire.NewSet(
	NewPlaceholder,
	wire.Bind(new(biz.RepositoryPlaceholder), new(*placeholder)),
)
