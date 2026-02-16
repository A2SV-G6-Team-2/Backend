package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DebtRepository interface {
	SumByDateRangeAndType(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, debtType string) (float64, error)
}
