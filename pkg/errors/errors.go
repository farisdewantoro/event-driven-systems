package errors

import (
	"eventdrivensystem/internal/generated/api_models"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

// Define a custom error struct to hold error code, message, and HTTP status
type APIError struct {
	ErrCode    string `json:"err_code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

// Implement the error interface for APIError
func (e *APIError) Error() string {
	return e.Message
}

// Helper function to create new API errors
func NewAPIError(code string, status int, message string, args ...interface{}) *APIError {
	return &APIError{
		ErrCode:    code,
		Message:    fmt.Sprintf(message, args...),
		HTTPStatus: status,
	}
}

func NewHTTPError(c echo.Context, err error) error {

	if e, ok := err.(validator.ValidationErrors); ok {
		resp := api_models.ErrorAPIResponse{
			ErrCode: ErrBadRequest.ErrCode,
			Message: parseValidationError(e[0]),
		}
		return echo.NewHTTPError(http.StatusBadRequest, resp)
	}

	if apiErr, ok := err.(*APIError); ok {
		resp := api_models.ErrorAPIResponse{
			ErrCode: apiErr.ErrCode,
			Message: apiErr.Message,
		}
		return echo.NewHTTPError(apiErr.HTTPStatus, resp)
	}

	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

// Helper function to parse validation errors
func parseValidationError(e validator.FieldError) string {
	var message string
	switch e.Tag() {
	case "required":
		message = fmt.Sprintf("The field '%s' is required.", e.Field())
	case "min":
		message = fmt.Sprintf("The field '%s' must be at least %s characters long.", e.Field(), e.Param())
	case "max":
		message = fmt.Sprintf("The field '%s' can be at most %s characters long.", e.Field(), e.Param())
	case "email":
		message = fmt.Sprintf("The field '%s' must be a valid email address.", e.Field())
	case "gte":
		message = fmt.Sprintf("The field '%s' must be greater than or equal to %s.", e.Field(), e.Param())
	case "lte":
		message = fmt.Sprintf("The field '%s' must be less than or equal to %s.", e.Field(), e.Param())
	default:
		message = fmt.Sprintf("Invalid value for the field '%s'.", e.Field())
	}
	return message
}
