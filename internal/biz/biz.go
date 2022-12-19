package biz

import "github.com/google/wire"

// provider set
var BizProviderSet = wire.NewSet(NewHealthzUsecase, NewUserUsecase)
