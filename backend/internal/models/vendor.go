package models

type Vendor struct {
	ID            int     `json:"id"`
	Name          string  `json:"name" validate:"required,min=2,max=255"`
	ContactPerson *string `json:"contact_person,omitempty"`
	Email         *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone         *string `json:"phone,omitempty"`
	Address       *string `json:"address,omitempty"`
}
