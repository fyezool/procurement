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
	ApproveRequisition(requisitionID int, adminID int) error
	RejectRequisition(requisitionID int, adminID int) error
	UpdateRequisition(requisitionID int, requesterID int, payload models.CreateRequisitionPayload) (*models.Requisition, error)
	DeleteRequisition(requisitionID int, requesterID int) error
	AdminUpdateRequisition(requisitionID int, adminID int, payload models.CreateRequisitionPayload) (*models.Requisition, error)
	AdminDeleteRequisition(requisitionID int, adminID int) error
}

type requisitionService struct {
	repo       repository.RequisitionRepository
	poService  PurchaseOrderService
	logService ActivityLogService
}

func NewRequisitionService(repo repository.RequisitionRepository, poService PurchaseOrderService, logService ActivityLogService) RequisitionService {
	return &requisitionService{repo: repo, poService: poService, logService: logService}
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

	createdReq, err := s.repo.CreateRequisition(requisition)
	if err != nil {
		details := err.Error()
		s.logService.Log(&requesterID, "CREATE_REQUISITION_FAILED", nil, nil, "FAILED", &details)
		return nil, err
	}

	s.logService.Log(&requesterID, "CREATE_REQUISITION_SUCCESS", Ptr("requisition"), &createdReq.ID, "SUCCESS", nil)
	return createdReq, nil
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

func (s *requisitionService) ApproveRequisition(requisitionID int, adminID int) error {
	// In a real app, we would wrap this in a transaction
	err := s.repo.UpdateRequisitionStatus(requisitionID, "Approved")
	if err != nil {
		details := err.Error()
		s.logService.Log(&adminID, "APPROVE_REQUISITION_FAILED", Ptr("requisition"), &requisitionID, "FAILED", &details)
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
		details := err.Error()
		s.logService.Log(&adminID, "APPROVE_REQUISITION_FAILED_PO_CREATION", Ptr("requisition"), &requisitionID, "FAILED", &details)
		return err
	}

	s.logService.Log(&adminID, "APPROVE_REQUISITION_SUCCESS", Ptr("requisition"), &requisitionID, "SUCCESS", nil)
	return nil
}

func (s *requisitionService) RejectRequisition(requisitionID int, adminID int) error {
	err := s.repo.UpdateRequisitionStatus(requisitionID, "Rejected")
	if err != nil {
		details := err.Error()
		s.logService.Log(&adminID, "REJECT_REQUISITION_FAILED", Ptr("requisition"), &requisitionID, "FAILED", &details)
		return err
	}
	s.logService.Log(&adminID, "REJECT_REQUISITION_SUCCESS", Ptr("requisition"), &requisitionID, "SUCCESS", nil)
	return nil
}

func (s *requisitionService) UpdateRequisition(requisitionID int, requesterID int, payload models.CreateRequisitionPayload) (*models.Requisition, error) {
	req, err := s.repo.GetRequisitionByID(requisitionID)
	if err != nil {
		return nil, err // Handles not found
	}

	if req.RequesterID != requesterID {
		return nil, ErrForbidden
	}

	if req.Status != "Pending" {
		return nil, ErrCannotModify
	}

	req.VendorID = payload.VendorID
	req.ItemDescription = payload.ItemDescription
	req.Quantity = payload.Quantity
	req.EstimatedPrice = payload.EstimatedPrice
	req.TotalPrice = payload.EstimatedPrice * float64(payload.Quantity)
	req.Justification = payload.Justification

	err = s.repo.UpdateRequisition(req)
	if err != nil {
		details := err.Error()
		s.logService.Log(&requesterID, "UPDATE_REQUISITION_FAILED", Ptr("requisition"), &requisitionID, "FAILED", &details)
		return nil, err
	}

	s.logService.Log(&requesterID, "UPDATE_REQUISITION_SUCCESS", Ptr("requisition"), &requisitionID, "SUCCESS", nil)
	return req, nil
}

func (s *requisitionService) DeleteRequisition(requisitionID int, requesterID int) error {
	req, err := s.repo.GetRequisitionByID(requisitionID)
	if err != nil {
		return err // Handles not found
	}

	if req.RequesterID != requesterID {
		return ErrForbidden
	}

	if req.Status != "Pending" {
		return ErrCannotModify
	}

	err = s.repo.DeleteRequisition(requisitionID)
	if err != nil {
		details := err.Error()
		s.logService.Log(&requesterID, "DELETE_REQUISITION_FAILED", Ptr("requisition"), &requisitionID, "FAILED", &details)
		return err
	}

	s.logService.Log(&requesterID, "DELETE_REQUISITION_SUCCESS", Ptr("requisition"), &requisitionID, "SUCCESS", nil)
	return nil
}

func (s *requisitionService) AdminUpdateRequisition(requisitionID int, adminID int, payload models.CreateRequisitionPayload) (*models.Requisition, error) {
	req, err := s.repo.GetRequisitionByID(requisitionID)
	if err != nil {
		return nil, err
	}

	// Admin can update any requisition, so no owner/status checks are needed.
	req.VendorID = payload.VendorID
	req.ItemDescription = payload.ItemDescription
	req.Quantity = payload.Quantity
	req.EstimatedPrice = payload.EstimatedPrice
	req.TotalPrice = payload.EstimatedPrice * float64(payload.Quantity)
	req.Justification = payload.Justification

	err = s.repo.UpdateRequisition(req)
	if err != nil {
		details := err.Error()
		s.logService.Log(&adminID, "ADMIN_UPDATE_REQUISITION_FAILED", Ptr("requisition"), &requisitionID, "FAILED", &details)
		return nil, err
	}

	s.logService.Log(&adminID, "ADMIN_UPDATE_REQUISITION_SUCCESS", Ptr("requisition"), &requisitionID, "SUCCESS", nil)
	return req, nil
}

func (s *requisitionService) AdminDeleteRequisition(requisitionID int, adminID int) error {
	// We might want to get the requisition first to log its details, but for now, this is fine.
	err := s.repo.DeleteRequisition(requisitionID)
	if err != nil {
		details := err.Error()
		s.logService.Log(&adminID, "ADMIN_DELETE_REQUISITION_FAILED", Ptr("requisition"), &requisitionID, "FAILED", &details)
		return err
	}
	s.logService.Log(&adminID, "ADMIN_DELETE_REQUISITION_SUCCESS", Ptr("requisition"), &requisitionID, "SUCCESS", nil)
	return nil
}
