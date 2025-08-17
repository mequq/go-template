package service

import (
	"context"
	"net/http"

	"github.com/google/wire"
)

//	@BasePath	/api/v1
var ServerProviderSet = wire.NewSet(NewHTTPHandler, http.NewServeMux)

// Service Interface
type Handler interface {
	RegisterHandler(ctx context.Context) error
}

//	@title			Swagger
//	@version		1.0
//	@description	API documentation
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://aban.io/support
//	@contact.email	info@aban.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
}
