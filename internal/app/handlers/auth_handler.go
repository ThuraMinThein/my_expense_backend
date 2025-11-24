package handlers

import (
	"net/http"

	"github.com/ThuraMinThein/my_expense_backend/config"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/api_structs"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/helper"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/services"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	service *services.AuthService
}

func (a *authHandler) SignUp(c *gin.Context) {
	var request api_structs.CreateUserRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_binding": err.Error()})
		return
	}

	user, err := a.service.SingUp(&request)
	if err != nil {
		if err.Error() == "username or email has already exist" {
			c.JSON(http.StatusConflict, gin.H{"error creating": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error_creating": err.Error()})
		return
	}

	setCookie(c, user.RefreshToken)

	c.JSON(http.StatusCreated, user)
}

func (a *authHandler) Login(c *gin.Context) {
	var request api_structs.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_binding": err.Error()})
		return
	}

	loginData, err := a.service.Login(&request)

	if err != nil {
		if err.Error() == "credential error" {
			c.JSON(http.StatusBadRequest, gin.H{"error_login": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error_login": err.Error()})
		return
	}

	setCookie(c, loginData.RefreshToken)

	c.JSON(http.StatusOK, loginData)
}

func (a *authHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claim, err := helper.ParseToken(refreshToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userToken, err := a.service.Refresh(claim.Sub, refreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setCookie(c, userToken.RefreshToken)

	c.JSON(http.StatusOK, userToken)
}

func (a *authHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claim, err := helper.ParseToken(refreshToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = a.service.Logout(claim.Sub, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	removeCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func setCookie(c *gin.Context, refreshToken string) {
	secure := config.Config.GinMode == "release"
	domain := config.Config.Domain

	// Set SameSite policy first
	if config.Config.GinMode == "release" {
		c.SetSameSite(http.SameSiteNoneMode) // For cross-site in production
	} else {
		c.SetSameSite(http.SameSiteLaxMode) // For local development
	}

	c.SetCookie("refreshToken", refreshToken, 7*24*60*60, "/", domain, secure, true)
}

func removeCookie(c *gin.Context) {
	domain := config.Config.Domain
	secure := config.Config.GinMode == "release"

	c.SetCookie("refreshToken", "", -1, "/", domain, secure, true)
}
