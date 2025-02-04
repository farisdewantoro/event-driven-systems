package rest

import (
	"eventdrivensystem/internal/usecase"

	goValidator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RouterHandler struct {
	echo      *echo.Echo
	validator *goValidator.Validate
	uc        *usecase.Usecase
}

func NewRouterHandler(echo *echo.Echo, validator *goValidator.Validate, uc *usecase.Usecase) *RouterHandler {
	return &RouterHandler{
		echo:      echo,
		validator: validator,
		uc:        uc,
	}
}

func (r *RouterHandler) RegisterRoutes() {
	base := r.echo.Group("/api")

	r.RegisterUserRoutes(base)
}
