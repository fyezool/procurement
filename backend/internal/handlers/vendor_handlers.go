package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
	"procurement-system/internal/services"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type VendorHandler struct {
	service  services.VendorService
	validate *validator.Validate
}

func NewVendorHandler(service services.VendorService) *VendorHandler {
	return &VendorHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *VendorHandler) CreateVendor(w http.ResponseWriter, r *http.Request) {
	var vendor models.Vendor
	if err := json.NewDecoder(r.Body).Decode(&vendor); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(&vendor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateVendor(&vendor); err != nil {
		http.Error(w, "Failed to create vendor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vendor)
}

func (h *VendorHandler) GetAllVendors(w http.ResponseWriter, r *http.Request) {
	vendors, err := h.service.GetAllVendors()
	if err != nil {
		http.Error(w, "Failed to retrieve vendors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendors)
}

func (h *VendorHandler) GetVendorByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid vendor ID", http.StatusBadRequest)
		return
	}

	vendor, err := h.service.GetVendorByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrVendorNotFound) {
			http.Error(w, "Vendor not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve vendor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

func (h *VendorHandler) UpdateVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid vendor ID", http.StatusBadRequest)
		return
	}

	var vendor models.Vendor
	if err := json.NewDecoder(r.Body).Decode(&vendor); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(&vendor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vendor.ID = id
	if err := h.service.UpdateVendor(&vendor); err != nil {
		if errors.Is(err, repository.ErrVendorNotFound) {
			http.Error(w, "Vendor not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update vendor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vendor)
}

func (h *VendorHandler) DeleteVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid vendor ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteVendor(id); err != nil {
		if errors.Is(err, repository.ErrVendorNotFound) {
			http.Error(w, "Vendor not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete vendor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
