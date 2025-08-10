package handlers

import (
	"encoding/json"
	"net/http"
	"procurement-system/internal/middleware"
	"strconv"

	"procurement-system/internal/models"
	"procurement-system/internal/repository"
	"procurement-system/internal/services"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// UserHandler handles HTTP requests for user management.
type UserHandler struct {
	userService services.UserService
	validate    *validator.Validate
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// GetAllUsers handles the request to retrieve all users.
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users: "+err.Error(), http.StatusInternalServerError)
	}
}

// GetUserByID handles the request to retrieve a single user by ID.
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user: "+err.Error(), http.StatusInternalServerError)
	}
}

// UpdateUser handles the request to update a user's details.
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	actorID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}

	var payload models.UpdateUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateUser(actorID, targetUserID, payload)
	if err != nil {
		if err == repository.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode updated user: "+err.Error(), http.StatusInternalServerError)
	}
}

// DeleteUser handles the request to delete a user.
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	actorID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}

	if err := h.userService.DeleteUser(actorID, targetUserID); err != nil {
		if err == repository.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
