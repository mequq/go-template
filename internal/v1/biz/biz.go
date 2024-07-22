package biz

import (
	biz "application/internal/v1/biz/healthz"
	"application/internal/v1/biz/sampleentity"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(sampleentity.NewSampleEntity, biz.NewHealthzBiz)
