package models

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"` // Do not expose password hash
	Role           string `json:"role"`
}

// RegistrationPayload defines the structure for user registration request
type RegistrationPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=Employee Admin 'Procurement Officer' Approver Vendor"`
}

// LoginPayload defines the structure for user login request
type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse defines the structure for a successful login response
type LoginResponse struct {
	Token string `json:"token"`
}
