package routes

import (
	"github.com/ThuraMinThein/my_expense_backend/internal/app/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *handlers.Handlers) {
	authRoutes(r, h)
	userRoutes(r, h)
	expenseRoutes(r, h)
}
