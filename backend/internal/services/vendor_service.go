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
	repo repository.VendorRepository
}

func NewVendorService(repo repository.VendorRepository) VendorService {
	return &vendorService{repo: repo}
}

func (s *vendorService) CreateVendor(vendor *models.Vendor) error {
	return s.repo.CreateVendor(vendor)
}

func (s *vendorService) GetAllVendors() ([]models.Vendor, error) {
	return s.repo.GetAllVendors()
}

func (s *vendorService) GetVendorByID(id int) (*models.Vendor, error) {
	return s.repo.GetVendorByID(id)
}

func (s *vendorService) UpdateVendor(vendor *models.Vendor) error {
	return s.repo.UpdateVendor(vendor)
}

func (s *vendorService) DeleteVendor(id int) error {
	return s.repo.DeleteVendor(id)
}
