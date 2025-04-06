package biz

import (
	"github.com/google/wire"
)

var BizProviderSet = wire.NewSet(NewSampleEntity, NewHealthzBiz)
