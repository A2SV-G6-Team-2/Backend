package http

import (
	"encoding/json"
	"net/http"

	"expense_tracker/usecases"

	"github.com/google/uuid"
)

type UserHandler struct {
	userUC usecases.UserUsecase
}

func NewUserHandler(uc usecases.UserUsecase) *UserHandler {
	return &UserHandler{uc}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.userUC.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var input usecases.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := h.userUC.Update(r.Context(), userID, input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
