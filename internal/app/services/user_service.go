package services

import (
	"mime/multipart"

	"github.com/ThuraMinThein/my_expense_backend/internal/app/api_structs"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/models"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/repositories"
)

type UserService struct {
	repository *repositories.Repositories
}

func (u *UserService) GetAll() ([]*models.User, error) {
	return u.repository.Users.GetAll()
}

func (u *UserService) GetOne(id uint64) (*models.User, error) {
	return u.repository.Users.GetOne(id)
}

func (u *UserService) Update(id uint64, profile_image *multipart.FileHeader, req *api_structs.UpdateUserRequest) (*models.User, error) {

	existingUser, err := u.GetOne(id)
	if err != nil {
		return nil, err
	}

	updatedUser := convertToModelUpdate(req)
	updatedUser.ID = existingUser.ID

	user, err := u.repository.Users.Update(updatedUser)

	return user, nil
}

// conversions
func convertToModel(user *api_structs.CreateUserRequest) *models.User {
	return &models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}

func convertToModelUpdate(user *api_structs.UpdateUserRequest) *models.User {
	return &models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}
