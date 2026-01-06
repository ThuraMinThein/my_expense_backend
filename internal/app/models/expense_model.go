package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expense struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	Amount      string         `gorm:"not null" json:"amount"` // Encrypted
	Category    string         `gorm:"not null" json:"category"` // Encrypted
	Note        string         `json:"note"`                   // Encrypted
	ExpenseDate time.Time      `gorm:"index" json:"expense_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Expense) TableName() string {
	return "expenses"
}
