package services

import (
	"errors"

	"github.com/ThuraMinThein/my_expense_backend/internal/app/api_structs"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/helper"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/models"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/repositories"
)

type AuthService struct {
	repositories *repositories.Repositories
}

func (as *AuthService) SingUp(request *api_structs.CreateUserRequest) (*models.UserToken, error) {
	hasUser, err := as.hasUsernameOrEmail(request)
	if err != nil {
		return nil, err
	}
	if hasUser {
		return nil, errors.New("username or email has already exist")
	}

	hashedPassword, err := helper.Hash(request.Password)
	if err != nil {
		return nil, err
	}
	request.Password = hashedPassword
	userModel := convertToModel(request)

	err = as.repositories.Users.Create(userModel)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := helper.GetTokens(userModel.ID)
	if err != nil {
		return nil, err
	}

	userToken := &models.UserToken{
		UserId:       userModel.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	err = as.repositories.Users.CreateToken(userToken)
	if err != nil {
		return nil, err
	}

	return userToken, nil
}

func (as *AuthService) Login(request *api_structs.LoginRequest) (*models.UserToken, error) {

	user, err := as.repositories.Users.GetByEmailOrUsername(request.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("credential error")
	}

	err = helper.VerifyHashed(user.Password, request.Password)
	if err != nil {
		return nil, errors.New("credential error")
	}

	accessToken, refreshToken, err := helper.GetTokens(user.ID)
	if err != nil {
		return nil, err
	}

	userToken := &models.UserToken{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	err = as.repositories.Users.UpdateToken(userToken)
	if err != nil {
		return nil, err
	}

	return userToken, nil
}

func (as *AuthService) Refresh(userId uint64, token string) (*models.UserToken, error) {
	user, err := as.repositories.Users.GetByRefreshToken(userId, token)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	accessToken, refreshToken, err := helper.GetTokens(user.ID)
	if err != nil {
		return nil, err
	}

	userToken := &models.UserToken{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	err = as.repositories.Users.UpdateToken(userToken)
	if err != nil {
		return nil, err
	}

	return userToken, nil
}

func (as *AuthService) Logout(userId uint64, refreshToken string) error {
	if userId == 0 || refreshToken == "" {
		return errors.New("invalid logout data")
	}

	if err := as.repositories.Users.DeleteRefreshToken(userId, refreshToken); err != nil {
		return err
	}

	return nil
}

// utils
func (as *AuthService) hasUsernameOrEmail(req *api_structs.CreateUserRequest) (bool, error) {
	if req.Email != "" {
		userByEmail, err := as.repositories.Users.GetByEmail(req.Email)
		if err != nil {
			return false, err
		}
		if userByEmail != nil {
			return true, nil
		}
	}

	userByUsername, err := as.repositories.Users.GetByUsername(req.Username)
	if err != nil {
		return false, err
	}
	if userByUsername != nil {
		return true, nil
	}

	return false, nil
}
