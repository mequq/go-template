package biz

import (
	biz "application/internal/biz/healthz"
	"application/internal/biz/sample_entity"
	"github.com/google/wire"
)

var BizProviderSet = wire.NewSet(sample_entity.NewSampleEntity, biz.NewHealthzBiz)
