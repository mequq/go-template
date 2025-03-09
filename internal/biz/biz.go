package biz

import (
	healthzusecase "application/internal/biz/healthz"
	sampleentityusecasev1 "application/internal/biz/sample"

	"github.com/google/wire"
)

var BizProviderSet = wire.NewSet(sampleentityusecasev1.NewSampleEntity, healthzusecase.NewHealthzBiz)
