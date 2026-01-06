package middlewares

import (
	"net/http"
	"strings"

	"github.com/ThuraMinThein/my_expense_backend/db"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/helper"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/repositories"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/services"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		tokenType := strings.Split(token, " ")

		if token == "" || len(tokenType) != 2 || tokenType[0] != "Bearer" {
			abortError(c, http.StatusUnauthorized)
			return
		}

		jwt := tokenType[1]

		claims, err := helper.ParseToken(jwt)

		if err != nil {
			abortError(c, http.StatusUnauthorized)
			return
		}

		repo := repositories.NewRepository(db.DB)
		services := services.NewServices(&repo)
		user, err := services.Users.GetOne(claims.Sub)
		if err != nil {
			abortError(c, http.StatusUnauthorized)
			return
		}
		c.Set("user", user)
		c.Set("user_id", user.ID)

		c.Next()

	}
}

func abortError(c *gin.Context, status int, message ...string) {
	errorMessage := ""
	switch status {
	case http.StatusUnauthorized:
		errorMessage = "Unauthorized"
	case http.StatusForbidden:
		errorMessage = "Forbidden"
	default:
		errorMessage = "Error"
	}
	if len(message) > 0 {
		errorMessage = errorMessage + ": " + message[0]
	}
	c.JSON(status, gin.H{"error": errorMessage})
	c.Abort()
}
