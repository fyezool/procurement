package services

import (
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
)

type VendorService interface {
	CreateVendor(vendor *models.Vendor) error
	GetAllVendors() ([]models.Vendor, error)
	GetVendorByID(id int) (*models.Vendor, error)
	UpdateVendor(vendor *models.Vendor) error
	DeleteVendor(id int) error
}

type vendorService struct {
	repo       repository.VendorRepository
	logService ActivityLogService
}

func NewVendorService(repo repository.VendorRepository, logService ActivityLogService) VendorService {
	return &vendorService{repo: repo, logService: logService}
}

func (s *vendorService) CreateVendor(vendor *models.Vendor) error {
	err := s.repo.CreateVendor(vendor)
	if err != nil {
		details := err.Error()
		s.logService.Log(nil, "CREATE_VENDOR_FAILED", Ptr("vendor"), nil, "FAILED", &details)
		return err
	}
	s.logService.Log(nil, "CREATE_VENDOR_SUCCESS", Ptr("vendor"), &vendor.ID, "SUCCESS", nil)
	return nil
}

func (s *vendorService) GetAllVendors() ([]models.Vendor, error) {
	return s.repo.GetAllVendors()
}

func (s *vendorService) GetVendorByID(id int) (*models.Vendor, error) {
	return s.repo.GetVendorByID(id)
}

func (s *vendorService) UpdateVendor(vendor *models.Vendor) error {
	err := s.repo.UpdateVendor(vendor)
	if err != nil {
		details := err.Error()
		s.logService.Log(nil, "UPDATE_VENDOR_FAILED", Ptr("vendor"), &vendor.ID, "FAILED", &details)
		return err
	}
	s.logService.Log(nil, "UPDATE_VENDOR_SUCCESS", Ptr("vendor"), &vendor.ID, "SUCCESS", nil)
	return nil
}

func (s *vendorService) DeleteVendor(id int) error {
	err := s.repo.DeleteVendor(id)
	if err != nil {
		details := err.Error()
		s.logService.Log(nil, "DELETE_VENDOR_FAILED", Ptr("vendor"), &id, "FAILED", &details)
		return err
	}
	s.logService.Log(nil, "DELETE_VENDOR_SUCCESS", Ptr("vendor"), &id, "SUCCESS", nil)
	return nil
}
