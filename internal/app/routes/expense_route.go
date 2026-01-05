package routes

import (
	"github.com/ThuraMinThein/my_expense_backend/internal/app/handlers"
	"github.com/ThuraMinThein/my_expense_backend/middlewares"
	"github.com/gin-gonic/gin"
)

func expenseRoutes(r *gin.Engine, h *handlers.Handlers) {
	protected := r.Group("/expenses").Use(middlewares.AuthMiddleware())
	{
		protected.POST("", h.ExpenseHandler.CreateExpense)
		protected.GET("", h.ExpenseHandler.GetExpenses)
		protected.DELETE("/:id", h.ExpenseHandler.DeleteExpense)
	}

	analytics := r.Group("/analytics").Use(middlewares.AuthMiddleware())
	{
		analytics.GET("/daily", h.ExpenseHandler.GetDailyUsage)
		analytics.GET("/weekly", h.ExpenseHandler.GetWeeklyUsage)
		analytics.GET("/monthly", h.ExpenseHandler.GetMonthlyUsage)
	}
}
