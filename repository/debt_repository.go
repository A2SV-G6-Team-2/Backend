package repository

import (
	"context"
	"expense_tracker/domain"
)

type DebtRepository interface {
	Create(ctx context.Context, debt *domain.Debt) error
	Update(ctx context.Context, debt *domain.Debt) error
	GetByID(ctx context.Context, id string) (*domain.Debt, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.Debt, error)
	ListUpcoming(ctx context.Context, userID string, days int) ([]*domain.Debt, error)
	MarkPaid(ctx context.Context, id string) (*domain.Debt, error)
	SetOverdue(ctx context.Context, nowUTC string) (int64, error)
	GetDueForReminder(ctx context.Context, nowUTC string) ([]*domain.Debt, error)
	UpdateReminder(ctx context.Context, id string, remindAtUTC string, sentAtUTC string) error
}