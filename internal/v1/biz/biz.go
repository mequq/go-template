package biz

import (
	biz "application/internal/v1/biz/healthz"
	"application/internal/v1/biz/sample_entity"
	"github.com/google/wire"
)

var BizProviderSet = wire.NewSet(sample_entity.NewSampleEntity, biz.NewHealthzBiz)
