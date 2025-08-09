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

// UpdateUserPayload defines the structure for updating a user's details.
// Admins can update a user's name and role. Email is not updatable for simplicity.
type UpdateUserPayload struct {
	Name string `json:"name" validate:"required"`
	Role string `json:"role" validate:"required,oneof=Employee Admin 'Procurement Officer' Approver Vendor"`
}

// UpdateProfilePayload defines the structure for updating a user's own name.
type UpdateProfilePayload struct {
	Name string `json:"name" validate:"required"`
}

// ChangePasswordPayload defines the structure for the change password request.
type ChangePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
