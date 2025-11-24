package services

import (
	"github.com/ThuraMinThein/my_expense_backend/internal/app/repositories"
)

type Services struct {
	Auth  *AuthService
	Users *UserService
}

func NewServices(repositories *repositories.Repositories) *Services {

	return &Services{
		Auth:  &AuthService{repositories},
		Users: &UserService{repository: repositories},
	}
}
