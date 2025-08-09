package models

import "time"

type Requisition struct {
	ID              int       `json:"id"`
	RequesterID     int       `json:"requester_id"`
	VendorID        *int      `json:"vendor_id"` // Pointer to allow null
	ItemDescription string    `json:"item_description" validate:"required"`
	Quantity        int       `json:"quantity" validate:"required,gt=0"`
	EstimatedPrice  float64   `json:"estimated_price" validate:"required,gt=0"`
	TotalPrice      float64   `json:"total_price"`
	Justification   string    `json:"justification"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateRequisitionPayload struct {
	VendorID        *int    `json:"vendor_id"`
	ItemDescription string  `json:"item_description" validate:"required"`
	Quantity        int     `json:"quantity" validate:"required,gt=0"`
	EstimatedPrice  float64 `json:"estimated_price" validate:"required,gt=0"`
	Justification   string  `json:"justification"`
}
