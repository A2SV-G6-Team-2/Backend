package http

import (
	"encoding/json"
	"expense_tracker/domain"
	"expense_tracker/usecases"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DebtHandler struct {
	usecase *usecases.DebtUsecase
}

func NewDebtHandler(usecase *usecases.DebtUsecase) *DebtHandler {
	return &DebtHandler{usecase: usecase}
}

type createDebtRequest struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	Type            string  `json:"type"`
	PeerName        string  `json:"peer_name"`
	Amount          float64 `json:"amount"`
	DueDate         string  `json:"due_date"`
	ReminderEnabled bool    `json:"reminder_enabled"`
	Note            *string `json:"note"`
}

type updateDebtRequest struct {
	Type            string  `json:"type"`
	PeerName        string  `json:"peer_name"`
	Amount          float64 `json:"amount"`
	DueDate         string  `json:"due_date"`
	ReminderEnabled bool    `json:"reminder_enabled"`
	Note            *string `json:"note"`
}

func (h *DebtHandler) CreateDebt(w http.ResponseWriter, r *http.Request) {
	var req createDebtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid due_date")
		return
	}

	debt := &domain.Debt{
		ID:              req.ID,
		UserID:          req.UserID,
		Type:            req.Type,
		PeerName:        req.PeerName,
		Amount:          req.Amount,
		DueDate:         dueDate,
		ReminderEnabled: req.ReminderEnabled,
		Note:            req.Note,
	}

	if err := h.usecase.Create(r.Context(), debt); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, debt)
}

func (h *DebtHandler) UpdateDebt(w http.ResponseWriter, r *http.Request) {
	debtID := extractDebtID(r.URL.Path)
	if debtID == "" {
		writeError(w, http.StatusBadRequest, "missing debt id")
		return
	}

	var req updateDebtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid due_date")
		return
	}

	debt := &domain.Debt{
		ID:              debtID,
		Type:            req.Type,
		PeerName:        req.PeerName,
		Amount:          req.Amount,
		DueDate:         dueDate,
		ReminderEnabled: req.ReminderEnabled,
		Note:            req.Note,
	}

	if err := h.usecase.Update(r.Context(), debt); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, debt)
}

func (h *DebtHandler) MarkDebtPaid(w http.ResponseWriter, r *http.Request) {
	debtID := extractDebtID(r.URL.Path)
	if debtID == "" {
		writeError(w, http.StatusBadRequest, "missing debt id")
		return
	}

	debt, err := h.usecase.MarkPaid(r.Context(), debtID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, debt)
}

func (h *DebtHandler) ListDebts(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	debts, err := h.usecase.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, debts)
}

func (h *DebtHandler) ListUpcomingDebts(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "missing user_id")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		parsed, err := parseInt(daysStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid days")
			return
		}
		days = parsed
	}

	debts, err := h.usecase.ListUpcoming(r.Context(), userID, days)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, debts)
}

func extractDebtID(path string) string {
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return ""
	}
	if parts[len(parts)-1] == "pay" && len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return parts[len(parts)-1]
}

func parseInt(value string) (int, error) {
	return strconv.Atoi(value)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
