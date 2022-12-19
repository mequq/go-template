package service

import (
	"app/internal/biz"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// tService is the t service.
type UserService struct {
	uc     *biz.UserUsecase
	logger zerolog.Logger
}

// NewTService creates a new t service.
func NewUserService(uc *biz.UserUsecase, logger zerolog.Logger) *UserService {

	sv := &UserService{
		uc:     uc,
		logger: logger,
	}
	// sv.RegisterRoutes2(app.Router)
	return sv
}

// Get is the get handler.
func (s *UserService) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.logger.Info().Interface("keys", c.Request.Header).Msg("test")
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		user, err := s.uc.Get(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// List is the list handler.
func (s *UserService) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := s.uc.List(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// register routes
func (s *UserService) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/user/:id", s.Get())
	r.GET("/users", s.List())
}
