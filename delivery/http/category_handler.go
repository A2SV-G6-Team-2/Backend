package http

import (
	"encoding/json"
	"net/http"

	"expense_tracker/domain"
	"expense_tracker/usecases"
)

// CategoryHandler handles category HTTP endpoints
type CategoryHandler struct {
	categoryUC *usecases.CategoryUseCase
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryUC *usecases.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUC: categoryUC}
}

// CreateCategoryRequest is the JSON body for POST /categories
type CreateCategoryRequest struct {
	Name   string  `json:"name"`
	UserID *string `json:"user_id,omitempty"`
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := UserIDFromRequest(r)
	// For create: if body has user_id we use it for user-defined; else we can create global (userID nil) or user's (use context)
	var createUserID *string
	if userID != "" {
		createUserID = &userID
	}

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	// If request specifies user_id, use it (for user-defined category)
	if req.UserID != nil {
		if *req.UserID != "" && !isValidUUID(*req.UserID) {
			http.Error(w, "user_id must be a valid UUID", http.StatusBadRequest)
			return
		}
		createUserID = req.UserID
	}

	input := domain.CreateCategoryInput{Name: req.Name, UserID: createUserID}
	cat, err := h.categoryUC.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(cat)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var userID *string
	if uid := UserIDFromRequest(r); uid != "" {
		userID = &uid
	}
	list, err := h.categoryUC.List(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isValidUUID(id) {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}
	var userID *string
	if uid := UserIDFromRequest(r); uid != "" {
		userID = &uid
	}
	cat, err := h.categoryUC.GetByID(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cat == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cat)
}

// UpdateCategoryRequest for PUT /categories/:id
type UpdateCategoryRequest struct {
	Name *string `json:"name,omitempty"`
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isValidUUID(id) {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}
	userIDStr := UserIDFromRequest(r)
	var userID *string
	if userIDStr != "" {
		userID = &userIDStr
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	input := domain.UpdateCategoryInput{Name: req.Name}
	cat, err := h.categoryUC.Update(r.Context(), id, userID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cat == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cat)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isValidUUID(id) {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}
	var userID *string
	if uid := UserIDFromRequest(r); uid != "" {
		userID = &uid
	}
	err := h.categoryUC.Delete(r.Context(), id, userID)
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
