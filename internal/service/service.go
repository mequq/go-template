package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// service provider set
var ServiceProviderSet = wire.NewSet(NewService, NewHealthzService)

// Response struct
type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// service interface
type ServiceInterface interface {
	RegisterRoutes(*gin.RouterGroup)
}

// service provider struct
type Service struct {
	svc []ServiceInterface
}

// NewService creates a new service.
func NewService(healthzSvc *HealthzService) *Service {
	svc := &Service{}
	svc.append(healthzSvc)
	return svc
}

// append service
func (s *Service) append(svc ServiceInterface) {
	s.svc = append(s.svc, svc)
}

// RegisterRoutes registers the routes.
func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	for _, svc := range s.svc {
		svc.RegisterRoutes(router)
	}
}
