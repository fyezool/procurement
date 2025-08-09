package services

import (
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
)

// UserService defines the interface for user management operations.
type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUser(id int, payload models.UpdateUserPayload) (*models.User, error)
	DeleteUser(id int) error
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// GetAllUsers retrieves all users.
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

// GetUserByID retrieves a single user by their ID.
func (s *userService) GetUserByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	// Do not expose password hash
	user.HashedPassword = ""
	return user, nil
}

// UpdateUser updates a user's details based on the provided payload.
func (s *userService) UpdateUser(id int, payload models.UpdateUserPayload) (*models.User, error) {
	// Get the existing user to ensure they exist before updating.
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err // Handles repository.ErrUserNotFound
	}

	// Update fields from payload.
	user.Name = payload.Name
	user.Role = payload.Role

	// Persist changes to the database.
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	// Return the updated user model, ensuring password hash is not exposed.
	user.HashedPassword = ""
	return user, nil
}

// DeleteUser deletes a user by their ID.
func (s *userService) DeleteUser(id int) error {
	return s.userRepo.DeleteUser(id)
}
