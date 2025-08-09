package repository

import (
	"database/sql"
	"fmt"
	"procurement-system/internal/models"
	"time"
)

var (
	ErrPurchaseOrderNotFound = sql.ErrNoRows
)

type PurchaseOrderRepository interface {
	CreatePurchaseOrder(po *models.PurchaseOrder) error
	GetPurchaseOrderByID(id int) (*models.PurchaseOrder, error)
	GetNextPONumber() (string, error)
	GetPDFData(poID int) (*models.PDFData, error)
}

type postgresPurchaseOrderRepository struct {
	db *sql.DB
}

func NewPostgresPurchaseOrderRepository(db *sql.DB) PurchaseOrderRepository {
	return &postgresPurchaseOrderRepository{db: db}
}

func (r *postgresPurchaseOrderRepository) CreatePurchaseOrder(po *models.PurchaseOrder) error {
	query := `
		INSERT INTO purchase_orders (po_number, requisition_id, vendor_id, order_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(
		query,
		po.PONumber, po.RequisitionID, po.VendorID, po.OrderDate,
	).Scan(&po.ID, &po.CreatedAt)
	return err
}

func (r *postgresPurchaseOrderRepository) GetPurchaseOrderByID(id int) (*models.PurchaseOrder, error) {
	po := &models.PurchaseOrder{}
	query := `
		SELECT po.id, po.po_number, po.requisition_id, po.vendor_id, po.order_date, po.created_at
		FROM purchase_orders po
		WHERE po.id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&po.ID, &po.PONumber, &po.RequisitionID, &po.VendorID, &po.OrderDate, &po.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (r *postgresPurchaseOrderRepository) GetNextPONumber() (string, error) {
	var count int
	year := time.Now().Year()
	query := `SELECT count(*) FROM purchase_orders WHERE EXTRACT(YEAR FROM order_date) = $1`
	err := r.db.QueryRow(query, year).Scan(&count)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PO-%d-%04d", year, count+1), nil
}

func (r *postgresPurchaseOrderRepository) GetPDFData(poID int) (*models.PDFData, error) {
	pdfData := &models.PDFData{}
	item := models.PDFItem{}

	query := `
		SELECT
			po.po_number,
			po.order_date,
			v.name,
			v.address,
			v.phone,
			v.email,
			r.item_description,
			r.quantity,
			r.estimated_price,
			r.total_price
		FROM purchase_orders po
		JOIN vendors v ON po.vendor_id = v.id
		JOIN requisitions r ON po.requisition_id = r.id
		WHERE po.id = $1
	`

	var orderDate time.Time
	var totalPrice float64

	err := r.db.QueryRow(query, poID).Scan(
		&pdfData.ReceiptNo,
		&orderDate,
		&pdfData.CustomerName,
		&pdfData.CustomerAddress,
		&pdfData.CustomerPhone,
		&pdfData.CustomerEmail,
		&item.Desc,
		&item.Qty,
		&item.UPrice,
		&totalPrice,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPurchaseOrderNotFound
		}
		return nil, err
	}

	pdfData.ReceiptDate = orderDate.Format("02/01/2006")
	pdfData.Items = []models.PDFItem{item}

	// These fields can be populated from a config or company profile in a real app
	pdfData.CompanyName = "Procurement Corp"
	pdfData.CompanyAddress = "123 Tech Avenue, Silicon Valley, CA 94043"
	pdfData.CompanyEmail = "contact@procurementcorp.com"
	pdfData.CompanyPhones = "1-800-555-PROC"
	pdfData.RegNo = "202301000001"
	pdfData.TinNo = "TIN12345678"
	pdfData.MsicCode = "62010"
	pdfData.BankName = "Global Tech Bank"
	pdfData.BankAccount = "123-456-7890"
	pdfData.PaymentMethod = "NET 30"


	return pdfData, nil
}
