package biz

import (
	"github.com/google/wire"

	biz "application/internal/v1/biz/healthz"
	"application/internal/v1/biz/sampleEntity"
)

var ProviderSet = wire.NewSet(sampleEntity.NewSampleEntity, biz.NewHealthzBiz)
