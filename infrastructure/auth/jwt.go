package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JWTService struct {
	Secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{Secret: secret}
}

func (j JWTService) Generate(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(10 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
}

func (j JWTService) Validate(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return uuid.Parse(claims["user_id"].(string))
}
