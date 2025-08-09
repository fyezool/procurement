package repository

import (
	"database/sql"
	"procurement-system/internal/models"
)

var (
	ErrVendorNotFound = sql.ErrNoRows
)

type VendorRepository interface {
	CreateVendor(vendor *models.Vendor) error
	GetAllVendors() ([]models.Vendor, error)
	GetVendorByID(id int) (*models.Vendor, error)
	UpdateVendor(vendor *models.Vendor) error
	DeleteVendor(id int) error
}

type postgresVendorRepository struct {
	db *sql.DB
}

func NewPostgresVendorRepository(db *sql.DB) VendorRepository {
	return &postgresVendorRepository{db: db}
}

func (r *postgresVendorRepository) CreateVendor(vendor *models.Vendor) error {
	query := `
		INSERT INTO vendors (name, contact_person, email, phone, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRow(query, vendor.Name, vendor.ContactPerson, vendor.Email, vendor.Phone, vendor.Address).Scan(&vendor.ID)
	return err
}

func (r *postgresVendorRepository) GetAllVendors() ([]models.Vendor, error) {
	query := `SELECT id, name, contact_person, email, phone, address FROM vendors`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		var v models.Vendor
		if err := rows.Scan(&v.ID, &v.Name, &v.ContactPerson, &v.Email, &v.Phone, &v.Address); err != nil {
			return nil, err
		}
		vendors = append(vendors, v)
	}
	return vendors, nil
}

func (r *postgresVendorRepository) GetVendorByID(id int) (*models.Vendor, error) {
	vendor := &models.Vendor{}
	query := `SELECT id, name, contact_person, email, phone, address FROM vendors WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&vendor.ID, &vendor.Name, &vendor.ContactPerson, &vendor.Email, &vendor.Phone, &vendor.Address)
	if err != nil {
		return nil, err // This will be sql.ErrNoRows if not found
	}
	return vendor, nil
}

func (r *postgresVendorRepository) UpdateVendor(vendor *models.Vendor) error {
	query := `
		UPDATE vendors
		SET name = $1, contact_person = $2, email = $3, phone = $4, address = $5
		WHERE id = $6
	`
	result, err := r.db.Exec(query, vendor.Name, vendor.ContactPerson, vendor.Email, vendor.Phone, vendor.Address, vendor.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrVendorNotFound
	}

	return nil
}

func (r *postgresVendorRepository) DeleteVendor(id int) error {
	query := `DELETE FROM vendors WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrVendorNotFound
	}

	return nil
}
