package repositories

import (
	"gorm.io/gorm"
)

type Repositories struct {
	Users   *UserStore
	Expense ExpenseRepository
}

func NewRepository(db *gorm.DB) Repositories {
	return Repositories{
		Users:   &UserStore{db},
		Expense: NewExpenseRepository(db),
	}
}
