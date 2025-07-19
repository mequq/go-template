package handler

import (
	"application/internal/service"

	"github.com/google/wire"
)

//	@title			Swagger
//	@version		2.0
//	@description	Application.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

var HandlerProviderSet = wire.NewSet(

	NewServiceList,
	NewMuxHealthzHandler,
)

// New ServiceList
func NewServiceList(
	healthzSvc *HealthzHandler,

) []service.Handler {
	return []service.Handler{
		healthzSvc,
	}
}
