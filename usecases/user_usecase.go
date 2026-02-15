package usecases

import (
	"context"
	"expense_tracker/domain"

	"github.com/google/uuid"
)

type UserUsecase interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (domain.User, error)
}
