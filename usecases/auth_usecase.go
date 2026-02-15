package usecases

import (
	"context"
	"expense_tracker/domain"
	"expense_tracker/infrastructure/auth"
	"expense_tracker/repository"
)

type AuthUsecase interface {
	Register(ctx context.Context, input RegisterInput) (domain.User, error)
	Login(ctx context.Context, input LoginInput) (AuthResponse, error)
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	Token string
	User  domain.User
}

type authUsecase struct {
	userRepo repository.UserRepository
	hasher   auth.PasswordHasher
	jwt      auth.JWTService
}
