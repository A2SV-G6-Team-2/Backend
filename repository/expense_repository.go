package repository

import (
	"context"
	"expense_tracker/domain"
)

// ExpenseRepository defines persistence for expenses
type ExpenseRepository interface {
	Create(ctx context.Context, input domain.CreateExpenseInput) (*domain.Expense, error)
	GetByID(ctx context.Context, id, userID string) (*domain.Expense, error)
	List(ctx context.Context, filter domain.ExpenseFilter) ([]*domain.Expense, error)
	Update(ctx context.Context, id, userID string, input domain.UpdateExpenseInput) (*domain.Expense, error)
	Delete(ctx context.Context, id, userID string) error
}
