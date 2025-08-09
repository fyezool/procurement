package services

import (
	"bytes"
	"errors"
	"procurement-system/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRequisitionRepository is a mock type for the RequisitionRepository
type MockRequisitionRepository struct {
	mock.Mock
}

func (m *MockRequisitionRepository) CreateRequisition(requisition *models.Requisition) (*models.Requisition, error) {
	args := m.Called(requisition)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Requisition), args.Error(1)
}
func (m *MockRequisitionRepository) GetRequisitionsByRequesterID(requesterID int) ([]models.Requisition, error) {
	args := m.Called(requesterID)
	return args.Get(0).([]models.Requisition), args.Error(1)
}
func (m *MockRequisitionRepository) GetPendingRequisitions() ([]models.Requisition, error) {
	args := m.Called()
	return args.Get(0).([]models.Requisition), args.Error(1)
}
func (m *MockRequisitionRepository) GetAllRequisitions() ([]models.Requisition, error) {
	args := m.Called()
	return args.Get(0).([]models.Requisition), args.Error(1)
}
func (m *MockRequisitionRepository) GetRequisitionByID(id int) (*models.Requisition, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Requisition), args.Error(1)
}
func (m *MockRequisitionRepository) UpdateRequisitionStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}
func (m *MockRequisitionRepository) UpdateRequisition(req *models.Requisition) error {
	args := m.Called(req)
	return args.Error(0)
}
func (m *MockRequisitionRepository) DeleteRequisition(id int) error {
	args := m.Called(id)
	return args.Error(0)
}


// MockPurchaseOrderService is a mock type for the PurchaseOrderService
type MockPurchaseOrderService struct {
	mock.Mock
}

func (m *MockPurchaseOrderService) CreatePurchaseOrderFromRequisition(requisition *models.Requisition) (*models.PurchaseOrder, error) {
	args := m.Called(requisition)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PurchaseOrder), args.Error(1)
}

func (m *MockPurchaseOrderService) GetPurchaseOrderByID(id int) (*models.PurchaseOrder, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PurchaseOrder), args.Error(1)
}

func (m *MockPurchaseOrderService) GetAllPurchaseOrders() ([]models.PurchaseOrder, error) {
	args := m.Called()
	return args.Get(0).([]models.PurchaseOrder), args.Error(1)
}

func (m *MockPurchaseOrderService) GeneratePurchaseOrderPDF(poID int) (*bytes.Buffer, error) {
	args := m.Called(poID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bytes.Buffer), args.Error(1)
}

func TestRequisitionService(t *testing.T) {
	t.Run("CreateRequisition", func(t *testing.T) {
		mockReqRepo := new(MockRequisitionRepository)
		mockPoService := new(MockPurchaseOrderService)
		requisitionService := NewRequisitionService(mockReqRepo, mockPoService)
		payload := models.CreateRequisitionPayload{
			ItemDescription: "Test Item",
			Quantity:        10,
			EstimatedPrice:  100,
		}

		// We expect CreateRequisition to be called and we return a mock requisition.
		mockReqRepo.On("CreateRequisition", mock.AnythingOfType("*models.Requisition")).Return(&models.Requisition{
			Status:     "Pending",
			TotalPrice: 1000,
		}, nil).Once()

		req, err := requisitionService.CreateRequisition(payload, 1)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, "Pending", req.Status)
		assert.Equal(t, float64(1000), req.TotalPrice)
		mockReqRepo.AssertExpectations(t)
	})

	t.Run("ApproveRequisition", func(t *testing.T) {
		mockReqRepo := new(MockRequisitionRepository)
		mockPoService := new(MockPurchaseOrderService)
		requisitionService := NewRequisitionService(mockReqRepo, mockPoService)
		reqID := 1
		vendorID := 123
		mockRequisition := &models.Requisition{ID: reqID, VendorID: &vendorID}

		mockReqRepo.On("UpdateRequisitionStatus", reqID, "Approved").Return(nil).Once()
		mockReqRepo.On("GetRequisitionByID", reqID).Return(mockRequisition, nil).Once()
		mockPoService.On("CreatePurchaseOrderFromRequisition", mockRequisition).Return(&models.PurchaseOrder{}, nil).Once()

		err := requisitionService.ApproveRequisition(reqID)
		assert.NoError(t, err)
		mockReqRepo.AssertExpectations(t)
		mockPoService.AssertExpectations(t)
	})

	t.Run("ApproveRequisition - UpdateStatus Fails", func(t *testing.T) {
		mockReqRepo := new(MockRequisitionRepository)
		mockPoService := new(MockPurchaseOrderService)
		requisitionService := NewRequisitionService(mockReqRepo, mockPoService)
		reqID := 2
		expectedErr := errors.New("update failed")
		mockReqRepo.On("UpdateRequisitionStatus", reqID, "Approved").Return(expectedErr).Once()

		err := requisitionService.ApproveRequisition(reqID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		mockReqRepo.AssertExpectations(t)
		mockPoService.AssertExpectations(t) // This mock was not called, so this should be fine.
		mockReqRepo.AssertNotCalled(t, "GetRequisitionByID", mock.Anything)
		mockPoService.AssertNotCalled(t, "CreatePurchaseOrderFromRequisition", mock.Anything)
	})

	t.Run("RejectRequisition", func(t *testing.T) {
		mockReqRepo := new(MockRequisitionRepository)
		mockPoService := new(MockPurchaseOrderService)
		requisitionService := NewRequisitionService(mockReqRepo, mockPoService)
		reqID := 3
		mockReqRepo.On("UpdateRequisitionStatus", reqID, "Rejected").Return(nil).Once()
		err := requisitionService.RejectRequisition(reqID)
		assert.NoError(t, err)
		mockReqRepo.AssertExpectations(t)
	})
}
