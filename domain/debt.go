package domain

import "time"

type DebtStatus string

const (
	DebtStatusPending DebtStatus = "pending"
	DebtStatusPaid    DebtStatus = "paid"
	DebtStatusOverdue DebtStatus = "overdue"
)

type Debt struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	Type            string     `json:"type"`
	PeerName        string     `json:"peer_name"`
	Amount          float64    `json:"amount"`
	DueDate         time.Time  `json:"due_date"`
	ReminderEnabled bool       `json:"reminder_enabled"`
	RemindAt        *time.Time `json:"remind_at,omitempty"`
	SentAt          *time.Time `json:"sent_at,omitempty"`
	Status          DebtStatus `json:"status"`
	Note            *string    `json:"note,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}
