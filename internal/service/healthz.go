package service

import (
	"app/internal/biz"

	"github.com/gin-gonic/gin"
)

// HealthzService is the healthz service.
type HealthzService struct {
	uc *biz.HealthzUsecase
}

// healthz response
type HealthzResponse struct {
	Message string `json:"message"`
}

// NewHealthzService creates a new healthz service.
func NewHealthzService(uc *biz.HealthzUsecase) *HealthzService {
	return &HealthzService{
		uc: uc,
	}
}

// readyness is the readyness handler.
func (s *HealthzService) readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := s.uc.Readiness(c)
		resp := Response{
			Code:   200,
			Status: "OK",
			Msg:    "success",
		}

		if err != nil {
			resp.Code = 500
			resp.Status = "ERROR"
			resp.Msg = err.Error()
		} else {
			resp.Code = 200
			resp.Status = "OK"
			resp.Msg = "success"
			resp.Data = HealthzResponse{
				Message: "OK",
			}
		}

		c.JSON(resp.Code, resp)
	}
}

// liveness is the liveness handler.
func (s *HealthzService) liveness() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := s.uc.Liveness(c)
		resp := Response{
			Code:   200,
			Status: "OK",
			Msg:    "success",
		}

		if err != nil {
			resp.Code = 500
			resp.Status = "ERROR"
			resp.Msg = err.Error()
		} else {
			resp.Code = 200
			resp.Status = "OK"
			resp.Msg = "success"
			resp.Data = HealthzResponse{
				Message: "OK",
			}
		}

		c.JSON(resp.Code, resp)
	}
}

// RegisterRoutes registers the routes.
func (s *HealthzService) RegisterRoutes(router *gin.RouterGroup) {
	router = router.Group("/healthz")
	router.GET("/readiness", s.readiness())
	router.GET("/liveness", s.liveness())
}
