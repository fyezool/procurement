package services

import (
	"errors"
	"procurement-system/internal/models"
	"procurement-system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIncorrectPassword = errors.New("incorrect old password")
)

// UserService defines the interface for user management operations.
type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUser(actorID int, targetUserID int, payload models.UpdateUserPayload) (*models.User, error)
	DeleteUser(actorID int, targetUserID int) error
	UpdateMyProfile(userID int, payload models.UpdateProfilePayload) (*models.User, error)
	ChangeMyPassword(userID int, payload models.ChangePasswordPayload) error
}

type userService struct {
	userRepo   repository.UserRepository
	logService ActivityLogService
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo repository.UserRepository, logService ActivityLogService) UserService {
	return &userService{userRepo: userRepo, logService: logService}
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
func (s *userService) UpdateUser(actorID int, targetUserID int, payload models.UpdateUserPayload) (*models.User, error) {
	// Get the existing user to ensure they exist before updating.
	user, err := s.userRepo.GetUserByID(targetUserID)
	if err != nil {
		return nil, err // Handles repository.ErrUserNotFound
	}

	// Update fields from payload.
	user.Name = payload.Name
	user.Role = payload.Role

	// Persist changes to the database.
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		details := err.Error()
		s.logService.Log(&actorID, "UPDATE_USER_FAILED", Ptr("user"), &targetUserID, "FAILED", &details)
		return nil, err
	}

	s.logService.Log(&actorID, "UPDATE_USER_SUCCESS", Ptr("user"), &targetUserID, "SUCCESS", nil)

	// Return the updated user model, ensuring password hash is not exposed.
	user.HashedPassword = ""
	return user, nil
}

// DeleteUser deletes a user by their ID.
func (s *userService) DeleteUser(actorID int, targetUserID int) error {
	err := s.userRepo.DeleteUser(targetUserID)
	if err != nil {
		details := err.Error()
		s.logService.Log(&actorID, "DELETE_USER_FAILED", Ptr("user"), &targetUserID, "FAILED", &details)
		return err
	}
	s.logService.Log(&actorID, "DELETE_USER_SUCCESS", Ptr("user"), &targetUserID, "SUCCESS", nil)
	return nil
}

func (s *userService) UpdateMyProfile(userID int, payload models.UpdateProfilePayload) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	user.Name = payload.Name

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		details := err.Error()
		s.logService.Log(&userID, "UPDATE_PROFILE_FAILED", Ptr("user"), &userID, "FAILED", &details)
		return nil, err
	}

	s.logService.Log(&userID, "UPDATE_PROFILE_SUCCESS", Ptr("user"), &userID, "SUCCESS", nil)
	user.HashedPassword = ""
	return user, nil
}

func (s *userService) ChangeMyPassword(userID int, payload models.ChangePasswordPayload) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Check if the old password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(payload.OldPassword))
	if err != nil {
		details := "Incorrect old password provided"
		s.logService.Log(&userID, "CHANGE_PASSWORD_FAILED", Ptr("user"), &userID, "FAILED", &details)
		return ErrIncorrectPassword
	}

	// Hash the new password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		details := err.Error()
		s.logService.Log(&userID, "CHANGE_PASSWORD_FAILED", Ptr("user"), &userID, "FAILED", &details)
		return err
	}

	// Update the password in the repository
	err = s.userRepo.UpdatePassword(userID, string(newHashedPassword))
	if err != nil {
		details := err.Error()
		s.logService.Log(&userID, "CHANGE_PASSWORD_FAILED", Ptr("user"), &userID, "FAILED", &details)
		return err
	}

	s.logService.Log(&userID, "CHANGE_PASSWORD_SUCCESS", Ptr("user"), &userID, "SUCCESS", nil)
	return nil
}
