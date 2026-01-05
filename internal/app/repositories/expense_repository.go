package repositories

import (
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
	return expenses, err
}

func (r *expenseRepository) GetByID(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Where("id = ?", id).First(&expense).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepository) Delete(id uuid.UUID, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Expense{}).Error
}

func (r *expenseRepository) GetDailyUsage(userID uint, date string) (float64, error) {
	var total float64
	err := r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND DATE(expense_date) = ?", userID, date).
		Scan(&total).Error
	return total, err
}

func (r *expenseRepository) GetWeeklyUsage(userID uint, week string) ([]map[string]interface{}, float64, error) {
	var dailyUsage []map[string]interface{}

	err := r.db.Model(&models.Expense{}).
		Select("DATE(expense_date) as date, COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND EXTRACT(YEAR FROM expense_date) = ? AND EXTRACT(WEEK FROM expense_date) = ?", userID, week[0:4], week[6:8]).
		Group("DATE(expense_date)").
		Order("date").
		Scan(&dailyUsage).Error

	if err != nil {
		return nil, 0, err
	}

	var weekTotal float64
	err = r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND EXTRACT(YEAR FROM expense_date) = ? AND EXTRACT(WEEK FROM expense_date) = ?", userID, week[0:4], week[6:8]).
		Scan(&weekTotal).Error

	return dailyUsage, weekTotal, err
}

func (r *expenseRepository) GetMonthlyUsageByCategory(userID uint, month string) ([]map[string]interface{}, float64, error) {
	var categoryUsage []map[string]interface{}

	err := r.db.Model(&models.Expense{}).
		Select("category, COALESCE(SUM(amount), 0) as amount").
		Where("user_id = ? AND EXTRACT(YEAR FROM expense_date) = ? AND EXTRACT(MONTH FROM expense_date) = ?", userID, month[0:4], month[5:7]).
		Group("category").
		Order("amount DESC").
		Scan(&categoryUsage).Error

	if err != nil {
		return nil, 0, err
	}

	var monthTotal float64
	err = r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND EXTRACT(YEAR FROM expense_date) = ? AND EXTRACT(MONTH FROM expense_date) = ?", userID, month[0:4], month[5:7]).
		Scan(&monthTotal).Error

	return categoryUsage, monthTotal, err
}
