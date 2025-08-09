package services

import (
	"bytes"
	"errors"
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
	"time"
)

type PurchaseOrderService interface {
	CreatePurchaseOrderFromRequisition(requisition *models.Requisition) (*models.PurchaseOrder, error)
	GetPurchaseOrderByID(id int) (*models.PurchaseOrder, error)
	GetAllPurchaseOrders() ([]models.PurchaseOrder, error)
	GeneratePurchaseOrderPDF(poID int) (*bytes.Buffer, error)
}

type purchaseOrderService struct {
	poRepo     repository.PurchaseOrderRepository
	vendRepo   repository.VendorRepository
	pdfService PDFService
}

func NewPurchaseOrderService(poRepo repository.PurchaseOrderRepository, vendRepo repository.VendorRepository, pdfService PDFService) PurchaseOrderService {
	return &purchaseOrderService{
		poRepo:     poRepo,
		vendRepo:   vendRepo,
		pdfService: pdfService,
	}
}

func (s *purchaseOrderService) CreatePurchaseOrderFromRequisition(requisition *models.Requisition) (*models.PurchaseOrder, error) {
	if requisition.VendorID == nil {
		return nil, errors.New("cannot create purchase order without a vendor")
	}

	poNumber, err := s.poRepo.GetNextPONumber()
	if err != nil {
		return nil, err
	}

	po := &models.PurchaseOrder{
		PONumber:      poNumber,
		RequisitionID: requisition.ID,
		VendorID:      *requisition.VendorID,
		OrderDate:     time.Now(),
	}

	err = s.poRepo.CreatePurchaseOrder(po)
	if err != nil {
		return nil, err
	}

	return po, nil
}

func (s *purchaseOrderService) GetPurchaseOrderByID(id int) (*models.PurchaseOrder, error) {
	return s.poRepo.GetPurchaseOrderByID(id)
}

func (s *purchaseOrderService) GetAllPurchaseOrders() ([]models.PurchaseOrder, error) {
	return s.poRepo.GetAllPurchaseOrders()
}

func (s *purchaseOrderService) GeneratePurchaseOrderPDF(poID int) (*bytes.Buffer, error) {
	pdfData, err := s.poRepo.GetPDFData(poID)
	if err != nil {
		return nil, err
	}

	return s.pdfService.GeneratePurchaseOrderPDF(pdfData)
}
