package repo

import (
	"application/internal/biz"

	"github.com/google/wire"
)

var RepoProvider = wire.NewSet(
	NewHealthzDS,
	NewSampleEntity,
	wire.Bind(new(biz.HealthzRepoInterface), new(*HealthzDS)),
	wire.Bind(new(biz.SampleEntityRepoInterface), new(*sampleEntity)),
)
