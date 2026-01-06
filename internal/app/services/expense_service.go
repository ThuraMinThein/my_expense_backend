package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/ThuraMinThein/my_expense_backend/internal/app/models"
	"github.com/ThuraMinThein/my_expense_backend/internal/app/repositories"
	"github.com/google/uuid"
)

type ExpenseService interface {
	CreateExpense(req CreateExpenseRequest, userID uint) (*models.Expense, error)
	GetExpenses(userID uint, from, to string) ([]models.Expense, error)
	DeleteExpense(id string, userID uint) error
	GetDailyUsage(userID uint, date string) (map[string]interface{}, error)
	GetWeeklyUsage(userID uint, week string) (map[string]interface{}, error)
	GetMonthlyUsage(userID uint, month string) (map[string]interface{}, error)
}

type expenseService struct {
	expenseRepo repositories.ExpenseRepository
}

type CreateExpenseRequest struct {
	Name        string  `json:"name" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Category    string  `json:"category" binding:"required"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

func NewExpenseService(expenseRepo repositories.ExpenseRepository) ExpenseService {
	return &expenseService{expenseRepo: expenseRepo}
}

func (s *expenseService) CreateExpense(req CreateExpenseRequest, userID uint) (*models.Expense, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	if strings.TrimSpace(req.Category) == "" {
		return nil, fmt.Errorf("category is required")
	}

	expenseDate := time.Now()
	if req.ExpenseDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.ExpenseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expense_date format, expected YYYY-MM-DD")
		}
		expenseDate = parsedDate
	}

	expense := &models.Expense{
		UserID:      userID,
		Name:        strings.TrimSpace(req.Name),
		Amount:      fmt.Sprintf("%f", req.Amount),
		Category:    strings.TrimSpace(req.Category),
		Note:        strings.TrimSpace(req.Note),
		ExpenseDate: expenseDate,
	}

	err := s.expenseRepo.Create(expense)
	if err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *expenseService) GetExpenses(userID uint, from, to string) ([]models.Expense, error) {
	return s.expenseRepo.GetByUserID(userID, from, to)
}

func (s *expenseService) DeleteExpense(id string, userID uint) error {
	expenseID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid expense ID format")
	}

	existingExpense, err := s.expenseRepo.GetByID(expenseID)
	if err != nil {
		return fmt.Errorf("expense not found")
	}

	if existingExpense.UserID != userID {
		return fmt.Errorf("unauthorized to delete this expense")
	}

	return s.expenseRepo.Delete(expenseID, userID)
}

func (s *expenseService) GetDailyUsage(userID uint, date string) (map[string]interface{}, error) {
	queryDate := date
	if queryDate == "" {
		queryDate = time.Now().Format("2006-01-02")
	}

	_, err := time.Parse("2006-01-02", queryDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}

	total, err := s.expenseRepo.GetDailyUsage(userID, queryDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"date":  queryDate,
		"total": total,
	}, nil
}

func (s *expenseService) GetWeeklyUsage(userID uint, week string) (map[string]interface{}, error) {
	queryWeek := week
	if queryWeek == "" {
		year, weekNum := time.Now().ISOWeek()
		queryWeek = fmt.Sprintf("%d-W%02d", year, weekNum)
	}

	if len(queryWeek) != 8 || queryWeek[4:5] != "-W" {
		return nil, fmt.Errorf("invalid week format, expected YYYY-WWW")
	}

	dailyUsage, weekTotal, err := s.expenseRepo.GetWeeklyUsage(userID, queryWeek)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"week":  queryWeek,
		"daily": dailyUsage,
		"total": weekTotal,
	}, nil
}

func (s *expenseService) GetMonthlyUsage(userID uint, month string) (map[string]interface{}, error) {
	queryMonth := month
	if queryMonth == "" {
		queryMonth = time.Now().Format("2006-01")
	}

	if len(queryMonth) != 7 || queryMonth[4:5] != "-" {
		return nil, fmt.Errorf("invalid month format, expected YYYY-MM")
	}

	_, err := time.Parse("2006-01", queryMonth)
	if err != nil {
		return nil, fmt.Errorf("invalid month format")
	}

	categoryUsage, monthTotal, err := s.expenseRepo.GetMonthlyUsageByCategory(userID, queryMonth)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"month":       queryMonth,
		"total":       monthTotal,
		"by_category": categoryUsage,
	}, nil
}
