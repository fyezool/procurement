package repository

import (
	"database/sql"
	"procurement-system/internal/models"
)

var (
	ErrRequisitionNotFound = sql.ErrNoRows
)

type RequisitionRepository interface {
	CreateRequisition(requisition *models.Requisition) error
	GetRequisitionsByRequesterID(requesterID int) ([]models.Requisition, error)
	GetPendingRequisitions() ([]models.Requisition, error)
	GetRequisitionByID(id int) (*models.Requisition, error)
	UpdateRequisitionStatus(id int, status string) error
}

type postgresRequisitionRepository struct {
	db *sql.DB
}

func NewPostgresRequisitionRepository(db *sql.DB) RequisitionRepository {
	return &postgresRequisitionRepository{db: db}
}

func (r *postgresRequisitionRepository) CreateRequisition(req *models.Requisition) error {
	query := `
		INSERT INTO requisitions (requester_id, vendor_id, item_description, quantity, estimated_price, total_price, justification, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(
		query,
		req.RequesterID, req.VendorID, req.ItemDescription, req.Quantity,
		req.EstimatedPrice, req.TotalPrice, req.Justification, req.Status,
	).Scan(&req.ID, &req.CreatedAt)
	return err
}

func (r *postgresRequisitionRepository) GetRequisitionsByRequesterID(requesterID int) ([]models.Requisition, error) {
	query := `
		SELECT id, requester_id, vendor_id, item_description, quantity, estimated_price, total_price, justification, status, created_at
		FROM requisitions
		WHERE requester_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, requesterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRequisitions(rows)
}

func (r *postgresRequisitionRepository) GetPendingRequisitions() ([]models.Requisition, error) {
	query := `
		SELECT id, requester_id, vendor_id, item_description, quantity, estimated_price, total_price, justification, status, created_at
		FROM requisitions
		WHERE status = 'Pending'
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRequisitions(rows)
}

func (r *postgresRequisitionRepository) GetRequisitionByID(id int) (*models.Requisition, error) {
	req := &models.Requisition{}
	query := `
		SELECT id, requester_id, vendor_id, item_description, quantity, estimated_price, total_price, justification, status, created_at
		FROM requisitions
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&req.ID, &req.RequesterID, &req.VendorID, &req.ItemDescription, &req.Quantity,
		&req.EstimatedPrice, &req.TotalPrice, &req.Justification, &req.Status, &req.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *postgresRequisitionRepository) UpdateRequisitionStatus(id int, status string) error {
	query := `UPDATE requisitions SET status = $1 WHERE id = $2`
	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRequisitionNotFound
	}

	return nil
}


func scanRequisitions(rows *sql.Rows) ([]models.Requisition, error) {
	var requisitions []models.Requisition
	for rows.Next() {
		var req models.Requisition
		if err := rows.Scan(
			&req.ID, &req.RequesterID, &req.VendorID, &req.ItemDescription, &req.Quantity,
			&req.EstimatedPrice, &req.TotalPrice, &req.Justification, &req.Status, &req.CreatedAt,
		); err != nil {
			return nil, err
		}
		requisitions = append(requisitions, req)
	}
	return requisitions, nil
}
