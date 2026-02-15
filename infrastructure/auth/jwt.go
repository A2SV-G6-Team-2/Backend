package auth

import "github.com/google/uuid"

type JWTService interface {
	Generate(userID uuid.UUID) (string, error)
	Validate(token string) (uuid.UUID, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password string, hash string) error
}
