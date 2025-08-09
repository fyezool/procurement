package handlers

import (
	"encoding/json"
	"net/http"
	"procurement-system/internal/middleware"
	"procurement-system/internal/services"
)

type NavigationHandler struct {
	service services.NavigationService
}

func NewNavigationHandler(service services.NavigationService) *NavigationHandler {
	return &NavigationHandler{service: service}
}

func (h *NavigationHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(middleware.UserRoleKey).(string)
	if !ok {
		http.Error(w, "Could not get user role from context", http.StatusInternalServerError)
		return
	}

	menu := h.service.GetMenuForRole(role)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func (h *NavigationHandler) GetBreadcrumbs(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(middleware.UserRoleKey).(string)
	if !ok {
		http.Error(w, "Could not get user role from context", http.StatusInternalServerError)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path query parameter is required", http.StatusBadRequest)
		return
	}

	breadcrumbs := h.service.GetBreadcrumbsForPath(role, path)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breadcrumbs)
}

