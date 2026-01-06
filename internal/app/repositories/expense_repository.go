package repositories

import (
	"fmt"

	"github.com/ThuraMinThein/my_expense_backend/internal/app/helper"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	GetByUserID(userID uint, from, to string) ([]models.Expense, error)
	GetByID(id uuid.UUID) (*models.Expense, error)
	Delete(id uuid.UUID, userID uint) error
	GetDailyUsage(userID uint, date string) (float64, error)
	GetWeeklyUsage(userID uint, week string) ([]map[string]interface{}, float64, error)
	GetMonthlyUsageByCategory(userID uint, month string) ([]map[string]interface{}, float64, error)
}

type expenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(expense *models.Expense) error {
	var err error
	expense.Name, err = helper.Encrypt(expense.Name)
	if err != nil {
		return err
	}
	expense.Amount, err = helper.Encrypt(expense.Amount)
	if err != nil {
		return err
	}
	expense.Category, err = helper.Encrypt(expense.Category)
	if err != nil {
		return err
	}
	if expense.Note != "" {
		expense.Note, err = helper.Encrypt(expense.Note)
		if err != nil {
			return err
		}
	}
	return r.db.Create(expense).Error
}

func (r *expenseRepository) GetByUserID(userID uint, from, to string) ([]models.Expense, error) {
	var expenses []models.Expense
	query := r.db.Where("user_id = ?", userID)

	if from != "" && to != "" {
		query = query.Where("expense_date BETWEEN ? AND ?", from, to)
	} else {
		query = query.Where("expense_date >= ?", "NOW() - INTERVAL '30 days'")
	}

	err := query.Order("expense_date DESC").Find(&expenses).Error
	if err != nil {
		return nil, err
	}

	for i := range expenses {
		expenses[i].Name, _ = helper.Decrypt(expenses[i].Name)
		expenses[i].Amount, _ = helper.Decrypt(expenses[i].Amount)
		expenses[i].Category, _ = helper.Decrypt(expenses[i].Category)
		if expenses[i].Note != "" {
			expenses[i].Note, _ = helper.Decrypt(expenses[i].Note)
		}
	}
	return expenses, err
}

func (r *expenseRepository) GetByID(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Where("id = ?", id).First(&expense).Error
	if err != nil {
		return nil, err
	}

	expense.Name, _ = helper.Decrypt(expense.Name)
	expense.Amount, _ = helper.Decrypt(expense.Amount)
	expense.Category, _ = helper.Decrypt(expense.Category)
	if expense.Note != "" {
		expense.Note, _ = helper.Decrypt(expense.Note)
	}

	return &expense, nil
}

func (r *expenseRepository) Delete(id uuid.UUID, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Expense{}).Error
}

func (r *expenseRepository) GetDailyUsage(userID uint, date string) (float64, error) {
	var expenses []models.Expense
	err := r.db.Where("user_id = ? AND DATE(expense_date) = ?", userID, date).Find(&expenses).Error
	if err != nil {
		return 0, err
	}

	var total float64
	for _, e := range expenses {
		decryptedAmount, err := helper.Decrypt(e.Amount)
		if err != nil {
			continue
		}
		var amount float64
		fmt.Sscanf(decryptedAmount, "%f", &amount)
		total += amount
	}
	return total, nil
}

func (r *expenseRepository) GetWeeklyUsage(userID uint, week string) ([]map[string]interface{}, float64, error) {
	var expenses []models.Expense
	err := r.db.Where("user_id = ? AND EXTRACT(YEAR FROM expense_date) = ? AND EXTRACT(WEEK FROM expense_date) = ?", userID, week[0:4], week[6:8]).
		Order("expense_date").
		Find(&expenses).Error

	if err != nil {
		return nil, 0, err
	}

	dailyUsageMap := make(map[string]float64)
	var weekTotal float64

	for _, e := range expenses {
		decryptedAmount, err := helper.Decrypt(e.Amount)
		if err != nil {
			continue
		}
		var amount float64
		fmt.Sscanf(decryptedAmount, "%f", &amount)
		weekTotal += amount

		dateStr := e.ExpenseDate.Format("2006-01-02")
		dailyUsageMap[dateStr] += amount
	}

	var dailyUsage []map[string]interface{}
	for date, total := range dailyUsageMap {
		dailyUsage = append(dailyUsage, map[string]interface{}{
			"date":  date,
			"total": total,
		})
	}

	return dailyUsage, weekTotal, nil
}

func (r *expenseRepository) GetMonthlyUsageByCategory(userID uint, month string) ([]map[string]interface{}, float64, error) {
	var expenses []models.Expense
	err := r.db.Where("user_id = ? AND EXTRACT(YEAR FROM expense_date) = ? AND EXTRACT(MONTH FROM expense_date) = ?", userID, month[0:4], month[5:7]).
		Find(&expenses).Error

	if err != nil {
		return nil, 0, err
	}

	categoryUsageMap := make(map[string]float64)
	var monthTotal float64

	for _, e := range expenses {
		decryptedAmount, err := helper.Decrypt(e.Amount)
		if err != nil {
			continue
		}
		var amount float64
		fmt.Sscanf(decryptedAmount, "%f", &amount)
		monthTotal += amount

		decryptedCategory, err := helper.Decrypt(e.Category)
		if err != nil {
			continue
		}
		categoryUsageMap[decryptedCategory] += amount
	}

	var categoryUsage []map[string]interface{}
	for category, amount := range categoryUsageMap {
		categoryUsage = append(categoryUsage, map[string]interface{}{
			"category": category,
			"amount":   amount,
		})
	}

	return categoryUsage, monthTotal, nil
}
