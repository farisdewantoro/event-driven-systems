package rest

import (
	"eventdrivensystem/internal/generated/api_models"
	"eventdrivensystem/internal/handler/rest/mapper"
	"eventdrivensystem/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *RouterHandler) RegisterUserRoutes(base *echo.Group) {
	v1 := base.Group("/v1/users")
	{

		v1.POST("", r.CreateUser)
	}
}

func (r *RouterHandler) CreateUser(c echo.Context) error {
	var (
		req  api_models.CreateUserRequest
		resp api_models.PlainResponse
		err  error
	)

	if err := c.Bind(&req); err != nil {
		return errors.NewHTTPError(c, errors.ErrBindRequest)
	}

	err = r.validator.Struct(&req)
	if err != nil {
		return errors.NewHTTPError(c, err)
	}

	err = r.uc.User.CreateUser(c.Request().Context(), mapper.ToCreateUserParam(&req))
	if err != nil {
		return errors.NewHTTPError(c, err)
	}

	resp.Success = true
	resp.Message = "User created successfully"

	return c.JSON(http.StatusOK, resp)
}
