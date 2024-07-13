package http

import "github.com/google/wire"

// @title Boilerplate
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/mahdimehrabi
// @contact.email mahdimehrabi17@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1
var ServerProviderSet = wire.NewSet(NewHttpHandler)
