package api_structs

import "github.com/ThuraMinThein/my_expense_backend/internal/app/models"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `binding:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Email    string `json:"email" form:"email"`
	Password string `form:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"-" form:"password"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"-"`
	User         models.User `json:"user"`
}
