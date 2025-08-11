package biz

import (
	"github.com/google/wire"
)

var BizProviderSet = wire.NewSet(
	NewHealthzBiz,
	wire.Bind(new(HealthzUseCaseInterface), new(*HealthzBiz)),
)
