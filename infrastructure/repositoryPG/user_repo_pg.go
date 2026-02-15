package repositoryPG

import (
	"context"
	"database/sql"
	"expense_tracker/domain"

	"github.com/google/uuid"
)

type UserRepoPG struct {
	DB *sql.DB
}

func NewUserRepoPG(db *sql.DB) *UserRepoPG {
	return &UserRepoPG{DB: db}
}

func (r *UserRepoPG) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users
	(user_id, name, email, password_hash, budgeting_style, default_currency, default_currency)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		u.UserID,
		u.Name,
		u.Email,
		u.PasswordHash,
		u.BudgetingStyle,
		u.DefaultCurrency,
	)
	return err
}

func (r *UserRepoPG) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}

	query := `SELECT user_id, name, email, password_hash, budgeting_style, default_currency, created_at
	FROM users
	WHERE email=$1`

	err := r.DB.QueryRowContext(ctx, query, email).
		Scan(&u.UserID, &u.Name, &u.Email, &u.PasswordHash, &u.BudgetingStyle, &u.DefaultCurrency, &u.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepoPG) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT user_id, name, email, password_hash, budgeting_style, default_currency, created_at
	FROM users
	WHERE user_id=$1`

	err := r.DB.QueryRowContext(ctx, query, id).
		Scan(&u.UserID, &u.Name, &u.Email, &u.BudgetingStyle, &u.DefaultCurrency, &u.CreatedAt)

	return u, err
}

func (r *UserRepoPG) Update(ctx context.Context, userID uuid.UUID, in domain.UpdateUserInput) error {
	query := `UPDATE users
	SET
		name = $1,
		budgeting_style = $2,
		default_currency = $3
	WHERE user_id = $4`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		in.Name,
		in.BudgetingStyle,
		in.DefaultCurrency,
		userID,
	)
	return err
}
