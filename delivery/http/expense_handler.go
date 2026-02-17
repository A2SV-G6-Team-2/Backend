package http

import (
	"encoding/json"
	"net/http"
	"time"

	"expense_tracker/domain"
	"expense_tracker/usecases"

	"github.com/google/uuid"
)

// ExpenseHandler handles expense HTTP endpoints
type ExpenseHandler struct {
	expenseUC *usecases.ExpenseUseCase
}

// NewExpenseHandler creates a new expense handler
func NewExpenseHandler(expenseUC *usecases.ExpenseUseCase) *ExpenseHandler {
	return &ExpenseHandler{expenseUC: expenseUC}
}

// CreateExpenseRequest is the JSON body for POST /expenses
type CreateExpenseRequest struct {
	ID              string   `json:"id"`
	Amount          float64  `json:"amount"`
	CategoryID      *string  `json:"category_id,omitempty"`
	IsRecurring     bool     `json:"is_recurring"`
	RecurrenceType  string   `json:"recurrence_type,omitempty"`
	NextDueDate     *string  `json:"next_due_date,omitempty"` // YYYY-MM-DD
	ReminderEnabled bool     `json:"reminder_enabled"`
	Note            string   `json:"note,omitempty"`
	ExpenseDate     string   `json:"expense_date"` // YYYY-MM-DD required
}

// UpdateExpenseRequest is the JSON body for PUT /expenses/:id
type UpdateExpenseRequest struct {
	Amount          *float64 `json:"amount,omitempty"`
	CategoryID      *string  `json:"category_id,omitempty"`
	IsRecurring     *bool    `json:"is_recurring,omitempty"`
	RecurrenceType  *string  `json:"recurrence_type,omitempty"`
	NextDueDate     *string  `json:"next_due_date,omitempty"`
	ReminderEnabled *bool    `json:"reminder_enabled,omitempty"`
	Note            *string  `json:"note,omitempty"`
	ExpenseDate     *string  `json:"expense_date,omitempty"`
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := UserIDFromRequest(r)
	if userID == "" {
		http.Error(w, "authorization required", http.StatusUnauthorized)
		return
	}

	var req CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Validation
	if req.Amount <= 0 {
		http.Error(w, "amount must be positive", http.StatusBadRequest)
		return
	}
	expenseDate, err := parseDate(req.ExpenseDate)
	if err != nil {
		http.Error(w, "expense_date required (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}
	if req.ID != "" && !isValidUUID(req.ID) {
		http.Error(w, "id must be a valid UUID", http.StatusBadRequest)
		return
	}
	if req.CategoryID != nil && *req.CategoryID != "" && !isValidUUID(*req.CategoryID) {
		http.Error(w, "category_id must be a valid UUID", http.StatusBadRequest)
		return
	}
	recType := domain.RecurrenceType(req.RecurrenceType)
	if req.IsRecurring && recType != domain.RecurrenceDaily && recType != domain.RecurrenceWeekly && recType != domain.RecurrenceMonthly {
		http.Error(w, "recurrence_type must be daily, weekly, or monthly when is_recurring is true", http.StatusBadRequest)
		return
	}
	var nextDue *time.Time
	if req.NextDueDate != nil && *req.NextDueDate != "" {
		t, err := parseDate(*req.NextDueDate)
		if err != nil {
			http.Error(w, "next_due_date invalid (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		nextDue = &t
	}

	expenseID := req.ID
	if expenseID == "" {
		expenseID = uuid.New().String()
	}

	input := domain.CreateExpenseInput{
		ID:              expenseID,
		UserID:          userID,
		Amount:          req.Amount,
		CategoryID:      req.CategoryID,
		IsRecurring:     req.IsRecurring,
		RecurrenceType:  recType,
		NextDueDate:     nextDue,
		ReminderEnabled: req.ReminderEnabled,
		Note:            req.Note,
		ExpenseDate:     expenseDate,
	}
	expense, err := h.expenseUC.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(expense)
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := UserIDFromRequest(r)
	if userID == "" {
		http.Error(w, "authorization required", http.StatusUnauthorized)
		return
	}

	var fromDate, toDate *time.Time
	if s := r.URL.Query().Get("from_date"); s != "" {
		t, err := parseDate(s)
		if err != nil {
			http.Error(w, "from_date invalid (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		fromDate = &t
	}
	if s := r.URL.Query().Get("to_date"); s != "" {
		t, err := parseDate(s)
		if err != nil {
			http.Error(w, "to_date invalid (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		toDate = &t
	}
	var categoryID *string
	if s := r.URL.Query().Get("category_id"); s != "" {
		if !isValidUUID(s) {
			http.Error(w, "category_id must be a valid UUID", http.StatusBadRequest)
			return
		}
		categoryID = &s
	}

	filter := usecases.ParseExpenseFilter(userID, fromDate, toDate, categoryID)
	list, err := h.expenseUC.List(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := UserIDFromRequest(r)
	if userID == "" {
		http.Error(w, "authorization required", http.StatusUnauthorized)
		return
	}
	if !isValidUUID(id) {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	expense, err := h.expenseUC.GetByID(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if expense == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(expense)
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := UserIDFromRequest(r)
	if userID == "" {
		http.Error(w, "authorization required", http.StatusUnauthorized)
		return
	}
	if !isValidUUID(id) {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var req UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	input := domain.UpdateExpenseInput{}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
			return
		}
		input.Amount = req.Amount
	}
	input.CategoryID = req.CategoryID
	if req.CategoryID != nil && *req.CategoryID != "" && !isValidUUID(*req.CategoryID) {
		http.Error(w, "category_id must be a valid UUID", http.StatusBadRequest)
		return
	}
	input.IsRecurring = req.IsRecurring
	if req.RecurrenceType != nil {
		rt := domain.RecurrenceType(*req.RecurrenceType)
		input.RecurrenceType = &rt
	}
	if req.NextDueDate != nil {
		t, err := parseDate(*req.NextDueDate)
		if err != nil {
			http.Error(w, "next_due_date invalid (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		input.NextDueDate = &t
	}
	input.ReminderEnabled = req.ReminderEnabled
	input.Note = req.Note
	if req.ExpenseDate != nil {
		t, err := parseDate(*req.ExpenseDate)
		if err != nil {
			http.Error(w, "expense_date invalid (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		input.ExpenseDate = &t
	}

	expense, err := h.expenseUC.Update(r.Context(), id, userID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if expense == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(expense)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := UserIDFromRequest(r)
	if userID == "" {
		http.Error(w, "authorization required", http.StatusUnauthorized)
		return
	}
	if !isValidUUID(id) {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	err := h.expenseUC.Delete(r.Context(), id, userID)
	if err != nil {
		if isErrNoRows(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
