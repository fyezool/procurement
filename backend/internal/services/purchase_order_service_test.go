package services

import (
	"bytes"
	"errors"
	"procurement-system/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPurchaseOrderRepository is a mock type for the PurchaseOrderRepository
type MockPurchaseOrderRepository struct {
	mock.Mock
}

func (m *MockPurchaseOrderRepository) CreatePurchaseOrder(po *models.PurchaseOrder) error {
	args := m.Called(po)
	return args.Error(0)
}
func (m *MockPurchaseOrderRepository) GetPurchaseOrderByID(id int) (*models.PurchaseOrder, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PurchaseOrder), args.Error(1)
}
func (m *MockPurchaseOrderRepository) GetNextPONumber() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
func (m *MockPurchaseOrderRepository) GetPDFData(poID int) (*models.PDFData, error) {
	args := m.Called(poID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PDFData), args.Error(1)
}

// MockPDFService is a mock type for the PDFService
type MockPDFService struct {
	mock.Mock
}

func (m *MockPDFService) GeneratePurchaseOrderPDF(data *models.PDFData) (*bytes.Buffer, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bytes.Buffer), args.Error(1)
}

func TestPurchaseOrderService(t *testing.T) {
	t.Run("CreatePurchaseOrderFromRequisition", func(t *testing.T) {
		mockPoRepo := new(MockPurchaseOrderRepository)
		poService := NewPurchaseOrderService(mockPoRepo, nil, nil)
		vendorID := 1
		requisition := &models.Requisition{
			ID:       1,
			VendorID: &vendorID,
		}
		mockPoRepo.On("GetNextPONumber").Return("PO-2023-0001", nil).Once()
		mockPoRepo.On("CreatePurchaseOrder", mock.AnythingOfType("*models.PurchaseOrder")).Return(nil).Once()

		po, err := poService.CreatePurchaseOrderFromRequisition(requisition)
		assert.NoError(t, err)
		assert.NotNil(t, po)
		assert.Equal(t, "PO-2023-0001", po.PONumber)
		mockPoRepo.AssertExpectations(t)
	})

	t.Run("CreatePurchaseOrderFromRequisition - No Vendor", func(t *testing.T) {
		mockPoRepo := new(MockPurchaseOrderRepository)
		poService := NewPurchaseOrderService(mockPoRepo, nil, nil)
		requisition := &models.Requisition{ID: 2} // No VendorID
		po, err := poService.CreatePurchaseOrderFromRequisition(requisition)
		assert.Error(t, err)
		assert.Nil(t, po)
		assert.Equal(t, "cannot create purchase order without a vendor", err.Error())
	})

	t.Run("GeneratePurchaseOrderPDF", func(t *testing.T) {
		mockPoRepo := new(MockPurchaseOrderRepository)
		mockPdfService := new(MockPDFService)
		poService := NewPurchaseOrderService(mockPoRepo, nil, mockPdfService)
		poID := 1
		pdfData := &models.PDFData{CompanyName: "Test Corp"}
		pdfBuffer := new(bytes.Buffer)
		pdfBuffer.WriteString("fake pdf content")

		mockPoRepo.On("GetPDFData", poID).Return(pdfData, nil).Once()
		mockPdfService.On("GeneratePurchaseOrderPDF", pdfData).Return(pdfBuffer, nil).Once()

		resultBuffer, err := poService.GeneratePurchaseOrderPDF(poID)
		assert.NoError(t, err)
		assert.NotNil(t, resultBuffer)
		assert.Equal(t, "fake pdf content", resultBuffer.String())
		mockPoRepo.AssertExpectations(t)
		mockPdfService.AssertExpectations(t)
	})

	t.Run("GeneratePurchaseOrderPDF - Repo Fails", func(t *testing.T) {
		mockPoRepo := new(MockPurchaseOrderRepository)
		mockPdfService := new(MockPDFService)
		poService := NewPurchaseOrderService(mockPoRepo, nil, mockPdfService)
		poID := 2
		expectedErr := errors.New("db error")

		mockPoRepo.On("GetPDFData", poID).Return(nil, expectedErr).Once()

		_, err := poService.GeneratePurchaseOrderPDF(poID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockPoRepo.AssertExpectations(t)
		mockPdfService.AssertNotCalled(t, "GeneratePurchaseOrderPDF", mock.Anything)
	})
}
