package http

import "github.com/google/wire"

var ServerProviderSet = wire.NewSet(NewHttpHandler)
