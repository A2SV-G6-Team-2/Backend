package repositoryPG

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type DebtRepoPG struct {
	DB *sql.DB
}

func NewDebtRepoPG(db *sql.DB) *DebtRepoPG {
	return &DebtRepoPG{DB: db}
}

func (r *DebtRepoPG) SumByDateRangeAndType(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, debtType string) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0)
	FROM debts
	WHERE user_id = $1 AND type = $2 AND due_date >= $3 AND due_date <= $4`

	var total sql.NullFloat64
	if err := r.DB.QueryRowContext(ctx, query, userID, debtType, startDate, endDate).Scan(&total); err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return total.Float64, nil
}
