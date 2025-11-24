package routes

import (
	"github.com/ThuraMinThein/my_expense_backend/internal/app/handlers"
	"github.com/ThuraMinThein/my_expense_backend/middlewares"
	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.Engine, h *handlers.Handlers) {
	user := r.Group("/users")
	user.Use(middlewares.AuthMiddleware())
	{
		user.GET("/me", h.UserHandler.GetLoginUser)
		user.GET("/:id", h.UserHandler.GetOne)
		user.PATCH("/:id", h.UserHandler.Update)
		user.DELETE("/:id", h.UserHandler.Delete)
	}
}
