package services

import (
	"procurement-system/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVendorRepository is a mock type for the VendorRepository
type MockVendorRepository struct {
	mock.Mock
}

func (m *MockVendorRepository) CreateVendor(vendor *models.Vendor) error {
	args := m.Called(vendor)
	return args.Error(0)
}

func (m *MockVendorRepository) GetAllVendors() ([]models.Vendor, error) {
	args := m.Called()
	return args.Get(0).([]models.Vendor), args.Error(1)
}

func (m *MockVendorRepository) GetVendorByID(id int) (*models.Vendor, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Vendor), args.Error(1)
}

func (m *MockVendorRepository) UpdateVendor(vendor *models.Vendor) error {
	args := m.Called(vendor)
	return args.Error(0)
}

func (m *MockVendorRepository) DeleteVendor(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestVendorService(t *testing.T) {
	mockRepo := new(MockVendorRepository)
	vendorService := NewVendorService(mockRepo)

	vendor := &models.Vendor{ID: 1, Name: "Test Vendor"}

	t.Run("CreateVendor", func(t *testing.T) {
		mockRepo.On("CreateVendor", vendor).Return(nil).Once()
		err := vendorService.CreateVendor(vendor)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetAllVendors", func(t *testing.T) {
		vendors := []models.Vendor{*vendor}
		mockRepo.On("GetAllVendors").Return(vendors, nil).Once()
		result, err := vendorService.GetAllVendors()
		assert.NoError(t, err)
		assert.Equal(t, vendors, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetVendorByID", func(t *testing.T) {
		mockRepo.On("GetVendorByID", 1).Return(vendor, nil).Once()
		result, err := vendorService.GetVendorByID(1)
		assert.NoError(t, err)
		assert.Equal(t, vendor, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateVendor", func(t *testing.T) {
		mockRepo.On("UpdateVendor", vendor).Return(nil).Once()
		err := vendorService.UpdateVendor(vendor)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteVendor", func(t *testing.T) {
		mockRepo.On("DeleteVendor", 1).Return(nil).Once()
		err := vendorService.DeleteVendor(1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
