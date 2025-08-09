package services

import (
	"errors"
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
)

var (
	ErrForbidden    = errors.New("user does not have permission to perform this action")
	ErrCannotModify = errors.New("requisition cannot be modified in its current state")
)

type RequisitionService interface {
	CreateRequisition(payload models.CreateRequisitionPayload, requesterID int) (*models.Requisition, error)
	GetMyRequisitions(requesterID int) ([]models.Requisition, error)
	GetPendingRequisitions() ([]models.Requisition, error)
	GetAllRequisitions() ([]models.Requisition, error)
	ApproveRequisition(requisitionID int) error
	RejectRequisition(requisitionID int) error
	UpdateRequisition(requisitionID int, requesterID int, payload models.CreateRequisitionPayload) (*models.Requisition, error)
	DeleteRequisition(requisitionID int, requesterID int) error
}

type requisitionService struct {
	repo      repository.RequisitionRepository
	poService PurchaseOrderService
}

func NewRequisitionService(repo repository.RequisitionRepository, poService PurchaseOrderService) RequisitionService {
	return &requisitionService{repo: repo, poService: poService}
}

func (s *requisitionService) CreateRequisition(payload models.CreateRequisitionPayload, requesterID int) (*models.Requisition, error) {
	requisition := &models.Requisition{
		RequesterID:     requesterID,
		VendorID:        payload.VendorID,
		ItemDescription: payload.ItemDescription,
		Quantity:        payload.Quantity,
		EstimatedPrice:  payload.EstimatedPrice,
		TotalPrice:      payload.EstimatedPrice * float64(payload.Quantity),
		Justification:   payload.Justification,
		Status:          "Pending",
	}

	return s.repo.CreateRequisition(requisition)
}

func (s *requisitionService) GetMyRequisitions(requesterID int) ([]models.Requisition, error) {
	return s.repo.GetRequisitionsByRequesterID(requesterID)
}

func (s *requisitionService) GetPendingRequisitions() ([]models.Requisition, error) {
	return s.repo.GetPendingRequisitions()
}

func (s *requisitionService) GetAllRequisitions() ([]models.Requisition, error) {
	return s.repo.GetAllRequisitions()
}

func (s *requisitionService) ApproveRequisition(requisitionID int) error {
	// In a real app, we would wrap this in a transaction
	err := s.repo.UpdateRequisitionStatus(requisitionID, "Approved")
	if err != nil {
		return err
	}

	req, err := s.repo.GetRequisitionByID(requisitionID)
	if err != nil {
		// Here we might want to roll back the status update
		return err
	}

	_, err = s.poService.CreatePurchaseOrderFromRequisition(req)
	if err != nil {
		// Here we might want to roll back the status update
		return err
	}
	return nil
}

func (s *requisitionService) RejectRequisition(requisitionID int) error {
	return s.repo.UpdateRequisitionStatus(requisitionID, "Rejected")
}

func (s *requisitionService) UpdateRequisition(requisitionID int, requesterID int, payload models.CreateRequisitionPayload) (*models.Requisition, error) {
	// Get the original requisition
	req, err := s.repo.GetRequisitionByID(requisitionID)
	if err != nil {
		return nil, err // Handles not found
	}

	// Check permissions
	if req.RequesterID != requesterID {
		return nil, ErrForbidden
	}

	// Check status
	if req.Status != "Pending" {
		return nil, ErrCannotModify
	}

	// Update fields from payload
	req.VendorID = payload.VendorID
	req.ItemDescription = payload.ItemDescription
	req.Quantity = payload.Quantity
	req.EstimatedPrice = payload.EstimatedPrice
	req.TotalPrice = payload.EstimatedPrice * float64(payload.Quantity)
	req.Justification = payload.Justification

	// Persist changes
	err = s.repo.UpdateRequisition(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *requisitionService) DeleteRequisition(requisitionID int, requesterID int) error {
	// Get the original requisition
	req, err := s.repo.GetRequisitionByID(requisitionID)
	if err != nil {
		return err // Handles not found
	}

	// Check permissions
	if req.RequesterID != requesterID {
		return ErrForbidden
	}

	// Check status
	if req.Status != "Pending" {
		return ErrCannotModify
	}

	// Delete the requisition
	return s.repo.DeleteRequisition(requisitionID)
}
