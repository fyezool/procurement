package services

import (
	"bytes"
	"fmt"
	"procurement-system/internal/models"

	"github.com/jung-kurt/gofpdf"
)

type PDFService interface {
	GeneratePurchaseOrderPDF(data *models.PDFData) (*bytes.Buffer, error)
}

type pdfService struct{}

func NewPDFService() PDFService {
	return &pdfService{}
}

func (s *pdfService) GeneratePurchaseOrderPDF(data *models.PDFData) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.Cell(40, 10, "Purchase Order")
	pdf.Ln(20)

	// Company and Vendor Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(100, 10, data.CompanyName)
	pdf.Cell(40, 10, "Vendor:")
	pdf.Ln(5)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(100, 10, data.CompanyAddress)
	pdf.Cell(40, 10, data.CustomerName)
	pdf.Ln(5)
	pdf.Cell(100, 10, data.CompanyEmail)
	pdf.Cell(40, 10, data.CustomerAddress)
	pdf.Ln(5)
	pdf.Cell(100, 10, data.CompanyPhones)
	pdf.Cell(40, 10, data.CustomerEmail)
	pdf.Ln(20)

	// PO Details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "PO Number:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, data.ReceiptNo)
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Date:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, data.ReceiptDate)
	pdf.Ln(20)

	// Items Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(100, 10, "Description")
	pdf.Cell(30, 10, "Quantity")
	pdf.Cell(30, 10, "Unit Price")
	pdf.Cell(30, 10, "Total")
	pdf.Ln(10)

	// Items Table Body
	pdf.SetFont("Arial", "", 10)
	var total float64
	for _, item := range data.Items {
		itemTotal := item.Qty * item.UPrice
		total += itemTotal
		pdf.Cell(100, 10, item.Desc)
		pdf.Cell(30, 10, fmt.Sprintf("%.2f", item.Qty))
		pdf.Cell(30, 10, fmt.Sprintf("%.2f", item.UPrice))
		pdf.Cell(30, 10, fmt.Sprintf("%.2f", itemTotal))
		pdf.Ln(5)
	}

	// Total
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(160, 10, "Total:")
	pdf.Cell(30, 10, fmt.Sprintf("%.2f", total))

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
