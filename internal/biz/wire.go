package biz

import (
	"github.com/google/wire"
)

var BizProviderSet = wire.NewSet(
	NewHealthz,
	wire.Bind(new(UsecaseHealthzer), new(*healthz)),

	NewPlaceholder,
	wire.Bind(new(UsecasePlaceholder), new(*placeholder)),
)
