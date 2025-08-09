package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"procurement-system/internal/repository"
	"procurement-system/internal/services"
	"strconv"

	"github.com/gorilla/mux"
)

type PurchaseOrderHandler struct {
	service services.PurchaseOrderService
}

func NewPurchaseOrderHandler(service services.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{service: service}
}

func (h *PurchaseOrderHandler) GetPurchaseOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid purchase order ID", http.StatusBadRequest)
		return
	}

	po, err := h.service.GetPurchaseOrderByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrPurchaseOrderNotFound) {
			http.Error(w, "Purchase order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve purchase order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(po)
}

func (h *PurchaseOrderHandler) GetPurchaseOrderPDF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid purchase order ID", http.StatusBadRequest)
		return
	}

	pdfBuffer, err := h.service.GeneratePurchaseOrderPDF(id)
	if err != nil {
		if errors.Is(err, repository.ErrPurchaseOrderNotFound) {
			http.Error(w, "Purchase order not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to generate PDF: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers to prompt download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"PO-%d.pdf\"", id))
	w.Header().Set("Content-Length", strconv.Itoa(pdfBuffer.Len()))

	// Write the PDF bytes to the response
	_, err = w.Write(pdfBuffer.Bytes())
	if err != nil {
		// Log the error, but the response has likely already been partially sent.
		http.Error(w, "Failed to write PDF to response", http.StatusInternalServerError)
	}
}
