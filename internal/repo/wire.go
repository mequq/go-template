package repo

import (
	"github.com/google/wire"
)

var RepoProvider = wire.NewSet(NewHealthzDS, NewSampleEntity, NewTokenRepo)
