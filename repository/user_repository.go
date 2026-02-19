package repository

import (
	"context"
	"expense_tracker/domain"

	"github.com/google/uuid"
)

// Implementation is in infrastructure/repository/user_repo_pg.go
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}
