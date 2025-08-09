package models

type Vendor struct {
	ID            int    `json:"id"`
	Name          string `json:"name" validate:"required"`
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email" validate:"required,email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
}
