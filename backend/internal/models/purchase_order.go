package models

import "time"

type PurchaseOrder struct {
	ID            int       `json:"id"`
	PONumber      string    `json:"po_number"`
	RequisitionID int       `json:"requisition_id"`
	VendorID      int       `json:"vendor_id"`
	OrderDate     time.Time `json:"order_date"`
	CreatedAt     time.Time `json:"created_at"`
}
