package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"expense_tracker/infrastructure/auth"
	"expense_tracker/usecases"
)

type UserHandler struct {
	userUC usecases.UserUsecase
	jwt    *auth.JWTService
}

func NewUserHandler(uc usecases.UserUsecase, jwt *auth.JWTService) *UserHandler {
	return &UserHandler{userUC: uc, jwt: jwt}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := h.jwt.Validate(tokenStr)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := h.jwt.Validate(tokenStr)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile updated successfully",
	})
}
