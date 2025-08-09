package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
	"procurement-system/internal/services"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService services.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload models.RegistrationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(payload)
	if err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(payload)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{Token: token})
}
