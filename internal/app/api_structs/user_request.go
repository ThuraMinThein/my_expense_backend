package api_structs

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `binding:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Email    string `json:"email" form:"email"`
	Password string `json:"-" form:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"-" form:"password"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}
