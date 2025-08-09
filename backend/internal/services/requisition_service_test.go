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

// We add this here to satisfy the interface for the mock.
func (m *MockPurchaseOrderService) AdminUpdateRequisition(requisitionID int, adminID int, payload models.CreateRequisitionPayload) (*models.Requisition, error) {
	return nil, nil
}
func (m *MockPurchaseOrderService) AdminDeleteRequisition(requisitionID int, adminID int) error {
	return nil
}


func TestRequisitionService(t *testing.T) {
	t.Run("CreateRequisition", func(t *testing.T) {
		mockReqRepo := new(MockRequisitionRepository)
		mockPoService := new(MockPurchaseOrderService)
		mockLogService := new(MockActivityLogService)
		requisitionService := NewRequisitionService(mockReqRepo, mockPoService, mockLogService)
		payload := models.CreateRequisitionPayload{
			ItemDescription: "Test Item",
			Quantity:        10,
			EstimatedPrice:  100,
		}

		mockReqRepo.On("CreateRequisition", mock.AnythingOfType("*models.Requisition")).Return(&models.Requisition{ID: 1, Status: "Pending", TotalPrice: 1000}, nil).Once()
		mockLogService.On("Log", mock.Anything, "CREATE_REQUISITION_SUCCESS", mock.Anything, mock.Anything, "SUCCESS", mock.Anything).Return()

		req, err := requisitionService.CreateRequisition(payload, 1)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		mockReqRepo.AssertExpectations(t)
		mockLogService.AssertExpectations(t)
	})

	t.Run("ApproveRequisition", func(t *testing.T) {
		mockReqRepo := new(MockRequisitionRepository)
		mockPoService := new(MockPurchaseOrderService)
		mockLogService := new(MockActivityLogService)
		requisitionService := NewRequisitionService(mockReqRepo, mockPoService, mockLogService)
		reqID := 1
		adminID := 99
		vendorID := 123
		mockRequisition := &models.Requisition{ID: reqID, VendorID: &vendorID}

		mockReqRepo.On("UpdateRequisitionStatus", reqID, "Approved").Return(nil).Once()
		mockReqRepo.On("GetRequisitionByID", reqID).Return(mockRequisition, nil).Once()
		mockPoService.On("CreatePurchaseOrderFromRequisition", mockRequisition).Return(&models.PurchaseOrder{}, nil).Once()
		mockLogService.On("Log", &adminID, "APPROVE_REQUISITION_SUCCESS", mock.Anything, &reqID, "SUCCESS", mock.Anything).Return()

		err := requisitionService.ApproveRequisition(reqID, adminID)
		assert.NoError(t, err)
		mockReqRepo.AssertExpectations(t)
		mockPoService.AssertExpectations(t)
		mockLogService.AssertExpectations(t)
	})

	t.Run("ApproveRequisition - UpdateStatus Fails", func(t *testing.T) {
		mockReqRepo := new(MockRequisitionRepository)
		mockPoService := new(MockPurchaseOrderService)
		mockLogService := new(MockActivityLogService)
		requisitionService := NewRequisitionService(mockReqRepo, mockPoService, mockLogService)
		reqID := 2
		adminID := 99
		expectedErr := errors.New("update failed")
		mockReqRepo.On("UpdateRequisitionStatus", reqID, "Approved").Return(expectedErr).Once()
		mockLogService.On("Log", &adminID, "APPROVE_REQUISITION_FAILED", mock.Anything, &reqID, "FAILED", mock.Anything).Return()

		err := requisitionService.ApproveRequisition(reqID, adminID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockReqRepo.AssertExpectations(t)
	})

	t.Run("RejectRequisition", func(t *testing.T) {
		mockReqRepo := new(MockRequisitionRepository)
		mockPoService := new(MockPurchaseOrderService)
		mockLogService := new(MockActivityLogService)
		requisitionService := NewRequisitionService(mockReqRepo, mockPoService, mockLogService)
		reqID := 3
		adminID := 99
		mockReqRepo.On("UpdateRequisitionStatus", reqID, "Rejected").Return(nil).Once()
		mockLogService.On("Log", &adminID, "REJECT_REQUISITION_SUCCESS", mock.Anything, &reqID, "SUCCESS", mock.Anything).Return()
		err := requisitionService.RejectRequisition(reqID, adminID)
		assert.NoError(t, err)
		mockReqRepo.AssertExpectations(t)
		mockLogService.AssertExpectations(t)
	})
}
