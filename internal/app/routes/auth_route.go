package routes

import (
	"github.com/ThuraMinThein/my_expense_backend/internal/app/handlers"
	"github.com/ThuraMinThein/my_expense_backend/middlewares"
	"github.com/gin-gonic/gin"
)

func authRoutes(r *gin.Engine, h *handlers.Handlers) {
	auth := r.Group("auth")
	{
		auth.POST("/sign-up", h.AuthHandler.SignUp)
		auth.POST("/login", h.AuthHandler.Login)
		auth.POST("/refresh", h.AuthHandler.Refresh)
	}
	protected := auth.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/logout", h.AuthHandler.Logout)
	}
}
