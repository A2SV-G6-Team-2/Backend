package usecases

import (
	"context"
	"errors"
	"expense_tracker/domain"
	"expense_tracker/repository"

	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password string, hash string) error
}

type JWTService interface {
	Generate(uuid.UUID) (string, error)
}

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
	hasher   PasswordHasher
	jwt      JWTService
}

func NewAuthUsecase(r repository.UserRepository, h PasswordHasher, j JWTService) AuthUsecase {
	return &authUsecase{r, h, j}
}

// Register and Login implementation
func (a *authUsecase) Register(ctx context.Context, in RegisterInput) (domain.User, error) {
	exists, _ := a.userRepo.GetByEmail(ctx, in.Email)
	if exists != nil {
		return domain.User{}, errors.New("email already used")
	}

	hash, err := a.hasher.Hash(in.Password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		UserID:          uuid.New(),
		Name:            in.Name,
		Email:           in.Email,
		PasswordHash:    hash,
		BudgetingStyle:  "flexible",
		DefaultCurrency: "ETB",
	}

	return user, a.userRepo.Create(ctx, &user)
}

func (a *authUsecase) Login(ctx context.Context, in LoginInput) (AuthResponse, error) {
	user, err := a.userRepo.GetByEmail(ctx, in.Email)
	if err != nil || user == nil {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	if err := a.hasher.Compare(in.Password, user.PasswordHash); err != nil {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	token, err := a.jwt.Generate(user.UserID)
	if err != nil {
		return AuthResponse{}, err
	}

	user.PasswordHash = ""
	return AuthResponse{Token: token, User: *user}, nil
}
