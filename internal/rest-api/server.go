package rest_api

import "github.com/google/wire"

var ServerProviderSet = wire.NewSet(NewHttpHandler)
