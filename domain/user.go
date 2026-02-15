package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID          uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	BudgetingStyle  string    `json:"budgeting_style"`
	DefaultCurrency string    `json:"default_currency"`
	CreatedAt       time.Time `json:"created_at"`
}
