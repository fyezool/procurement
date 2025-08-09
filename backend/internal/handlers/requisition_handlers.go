package handlers

import (
	"encoding/json"
	"net/http"
	"procurement-system/internal/middleware"
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
	"procurement-system/internal/services"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type RequisitionHandler struct {
	service  services.RequisitionService
	validate *validator.Validate
}

func NewRequisitionHandler(service services.RequisitionService) *RequisitionHandler {
	return &RequisitionHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *RequisitionHandler) CreateRequisition(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateRequisitionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requesterID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}

	requisition, err := h.service.CreateRequisition(payload, requesterID)
	if err != nil {
		http.Error(w, "Failed to create requisition", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(requisition)
}

func (h *RequisitionHandler) GetMyRequisitions(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}

	requisitions, err := h.service.GetMyRequisitions(requesterID)
	if err != nil {
		http.Error(w, "Failed to retrieve requisitions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requisitions)
}

func (h *RequisitionHandler) GetAllRequisitions(w http.ResponseWriter, r *http.Request) {
	requisitions, err := h.service.GetAllRequisitions()
	if err != nil {
		http.Error(w, "Failed to retrieve requisitions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requisitions)
}

func (h *RequisitionHandler) GetPendingRequisitions(w http.ResponseWriter, r *http.Request) {
	requisitions, err := h.service.GetPendingRequisitions()
	if err != nil {
		http.Error(w, "Failed to retrieve pending requisitions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requisitions)
}

func (h *RequisitionHandler) ApproveRequisition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid requisition ID", http.StatusBadRequest)
		return
	}

	if err := h.service.ApproveRequisition(id); err != nil {
		http.Error(w, "Failed to approve requisition", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Requisition approved successfully"})
}

func (h *RequisitionHandler) RejectRequisition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid requisition ID", http.StatusBadRequest)
		return
	}

	if err := h.service.RejectRequisition(id); err != nil {
		http.Error(w, "Failed to reject requisition", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Requisition rejected successfully"})
}

func (h *RequisitionHandler) UpdateRequisition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid requisition ID", http.StatusBadRequest)
		return
	}

	var payload models.CreateRequisitionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requesterID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}

	requisition, err := h.service.UpdateRequisition(id, requesterID, payload)
	if err != nil {
		switch err {
		case services.ErrForbidden:
			http.Error(w, err.Error(), http.StatusForbidden)
		case services.ErrCannotModify:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case repository.ErrRequisitionNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Failed to update requisition", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requisition)
}

func (h *RequisitionHandler) DeleteRequisition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid requisition ID", http.StatusBadRequest)
		return
	}

	requesterID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Could not get user ID from context", http.StatusInternalServerError)
		return
	}

	err = h.service.DeleteRequisition(id, requesterID)
	if err != nil {
		switch err {
		case services.ErrForbidden:
			http.Error(w, err.Error(), http.StatusForbidden)
		case services.ErrCannotModify:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case repository.ErrRequisitionNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Failed to delete requisition", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
