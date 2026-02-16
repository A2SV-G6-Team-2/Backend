package repositoryPG

import (
	"context"
	"database/sql"
	"expense_tracker/repository"
	"time"

	"github.com/google/uuid"
)

type ExpenseRepoPG struct {
	DB *sql.DB
}

func NewExpenseRepoPG(db *sql.DB) *ExpenseRepoPG {
	return &ExpenseRepoPG{DB: db}
}

func (r *ExpenseRepoPG) SumByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0)
	FROM expenses
	WHERE user_id = $1 AND expense_date >= $2 AND expense_date <= $3`

	var total sql.NullFloat64
	if err := r.DB.QueryRowContext(ctx, query, userID, startDate, endDate).Scan(&total); err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return total.Float64, nil
}

func (r *ExpenseRepoPG) CategoryBreakdownByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]repository.CategoryTotal, error) {
	query := `SELECT COALESCE(c.name, 'Uncategorized') AS category_name,
		COALESCE(SUM(e.amount), 0) AS total
	FROM expenses e
	LEFT JOIN categories c ON e.category_id = c.id
	WHERE e.user_id = $1 AND e.expense_date >= $2 AND e.expense_date <= $3
	GROUP BY category_name
	ORDER BY total DESC`

	rows, err := r.DB.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []repository.CategoryTotal{}
	for rows.Next() {
		var name string
		var total sql.NullFloat64
		if err := rows.Scan(&name, &total); err != nil {
			return nil, err
		}
		itemTotal := 0.0
		if total.Valid {
			itemTotal = total.Float64
		}
		results = append(results, repository.CategoryTotal{
			CategoryName: name,
			Total:        itemTotal,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
