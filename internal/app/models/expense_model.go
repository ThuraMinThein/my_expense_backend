package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expense struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uint      `gorm:"not null;index"`
	Amount      float64   `gorm:"not null"`
	Category    string    `gorm:"not null"`
	Note        string
	ExpenseDate time.Time `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Expense) TableName() string {
	return "expenses"
}
