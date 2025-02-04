package mapper

import (
	"eventdrivensystem/internal/generated/api_models"
	models "eventdrivensystem/internal/models/user"
)

func ToCreateUserParam(request *api_models.CreateUserRequest) *models.CreateUserParam {
	return &models.CreateUserParam{
		Email:    request.Email,
		Password: request.Password,
	}
}
