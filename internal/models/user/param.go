package models

type CreateUserParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (param *CreateUserParam) ToDomain() *User {
	return &User{
		Email:    param.Email,
		Password: param.Password,
	}
}
