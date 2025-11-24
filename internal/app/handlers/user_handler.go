package handlers

import (
	"net/http"
	"strconv"

	"github.com/ThuraMinThein/my_expense_backend/internal/app/api_structs"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/models"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/services"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	services *services.Services
}

func (U *userHandler) GetLoginUser(c *gin.Context) {
	userInterface, _ := c.Get("user")

	user, ok := userInterface.(*models.User)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *userHandler) GetAll(c *gin.Context) {
	users, err := u.services.Users.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)

}

func (u *userHandler) GetOne(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := u.services.Users.GetOne(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)

}

func (u *userHandler) Update(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var request api_structs.UpdateUserRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_binding": err.Error()})
		return
	}

	profile_image, _ := c.FormFile("profile_image")

	user, err := u.services.Users.Update(id, profile_image, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *userHandler) Delete(c *gin.Context) {

}
