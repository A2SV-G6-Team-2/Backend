package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CategoryTotal struct {
	CategoryName string
	Total        float64
}

type ExpenseRepository interface {
	SumByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (float64, error)
	CategoryBreakdownByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]CategoryTotal, error)
}
