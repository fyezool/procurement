package services

import (
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
)

type RequisitionService interface {
	CreateRequisition(payload models.CreateRequisitionPayload, requesterID int) (*models.Requisition, error)
	GetMyRequisitions(requesterID int) ([]models.Requisition, error)
	GetPendingRequisitions() ([]models.Requisition, error)
	ApproveRequisition(requisitionID int) error
	RejectRequisition(requisitionID int) error
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

	err := s.repo.CreateRequisition(requisition)
	if err != nil {
		return nil, err
	}

	return requisition, nil
}

func (s *requisitionService) GetMyRequisitions(requesterID int) ([]models.Requisition, error) {
	return s.repo.GetRequisitionsByRequesterID(requesterID)
}

func (s *requisitionService) GetPendingRequisitions() ([]models.Requisition, error) {
	return s.repo.GetPendingRequisitions()
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
