package handlers

import "github.com/ThuraMinThein/my_expense_backend/internal/app/services"

type Handlers struct {
	AuthHandler    *authHandler
	UserHandler    *userHandler
	ExpenseHandler *ExpenseHandler
}

func InitHandlers(services *services.Services) *Handlers {
	return &Handlers{
		AuthHandler:    &authHandler{service: services.Auth},
		UserHandler:    &userHandler{services: services},
		ExpenseHandler: NewExpenseHandler(services.Expense),
	}
}
