package services

import (
	"errors"
	"os"
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock type for the UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// These methods were added to the interface but are not used in this test file.
// We add them here to satisfy the interface.
func (m *MockUserRepository) GetUserByID(id int) (*models.User, error) { return nil, nil }
func (m *MockUserRepository) GetAllUsers() ([]models.User, error)    { return nil, nil }
func (m *MockUserRepository) UpdateUser(user *models.User) error      { return nil }
func (m *MockUserRepository) DeleteUser(id int) error                 { return nil }
func (m *MockUserRepository) UpdatePassword(userID int, newHashedPassword string) error { return nil }

// MockActivityLogService is a mock type for the ActivityLogService
type MockActivityLogService struct {
	mock.Mock
}

func (m *MockActivityLogService) Log(userID *int, action string, targetType *string, targetID *int, status string, details *string) {
	m.Called(userID, action, targetType, targetID, status, details)
}
func (m *MockActivityLogService) GetAll() ([]models.ActivityLog, error) {
	args := m.Called()
	return args.Get(0).([]models.ActivityLog), args.Error(1)
}


func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogService := new(MockActivityLogService)
	authService := NewAuthService(mockRepo, mockLogService)

	payload := models.RegistrationPayload{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "Employee",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(&models.User{ID: 1}, nil)
	mockLogService.On("Log", mock.Anything, "REGISTER_USER_SUCCESS", mock.Anything, mock.Anything, "SUCCESS", mock.Anything).Return()

	user, err := authService.Register(payload)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	mockRepo.AssertExpectations(t)
	mockLogService.AssertExpectations(t)
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogService := new(MockActivityLogService)
	authService := NewAuthService(mockRepo, mockLogService)
	payload := models.RegistrationPayload{Email: "exists@example.com"}

	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil, repository.ErrEmailExists)
	mockLogService.On("Log", mock.Anything, "REGISTER_USER_FAILED", mock.Anything, mock.Anything, "FAILED", mock.Anything).Return()

	user, err := authService.Register(payload)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, repository.ErrEmailExists, err)
	mockRepo.AssertExpectations(t)
	mockLogService.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogService := new(MockActivityLogService)
	authService := NewAuthService(mockRepo, mockLogService)
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUser := &models.User{ID: 1, Email: "test@example.com", HashedPassword: string(hashedPassword)}

	mockRepo.On("GetUserByEmail", "test@example.com").Return(mockUser, nil)
	mockLogService.On("Log", &mockUser.ID, "LOGIN_SUCCESS", mock.Anything, &mockUser.ID, "SUCCESS", mock.Anything).Return()

	token, err := authService.Login(models.LoginPayload{Email: "test@example.com", Password: "password123"})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
	mockLogService.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogService := new(MockActivityLogService)
	authService := NewAuthService(mockRepo, mockLogService)

	// Test case 1: User not found
	mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, repository.ErrUserNotFound)
	mockLogService.On("Log", mock.Anything, "LOGIN_FAILED", mock.Anything, mock.Anything, "FAILED", mock.Anything).Return().Once()
	_, err := authService.Login(models.LoginPayload{Email: "notfound@example.com", Password: "password"})
	assert.Equal(t, ErrInvalidCredentials, err)

	// Test case 2: Wrong password
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUser := &models.User{ID: 1, Email: "test@example.com", HashedPassword: string(hashedPassword)}
	mockRepo.On("GetUserByEmail", "test@example.com").Return(mockUser, nil)
	mockLogService.On("Log", &mockUser.ID, "LOGIN_FAILED", mock.Anything, &mockUser.ID, "FAILED", mock.Anything).Once()
	_, err = authService.Login(models.LoginPayload{Email: "test@example.com", Password: "wrongpassword"})
	assert.Equal(t, ErrInvalidCredentials, err)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogService := new(MockActivityLogService)
	authService := NewAuthService(mockRepo, mockLogService)
	expectedErr := errors.New("database error")
	mockRepo.On("GetUserByEmail", "any@example.com").Return(nil, expectedErr)
	mockLogService.On("Log", mock.Anything, "LOGIN_FAILED_DB_ERROR", mock.Anything, mock.Anything, "FAILED", mock.Anything).Return()

	_, err := authService.Login(models.LoginPayload{Email: "any@example.com", Password: "password"})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
	mockLogService.AssertExpectations(t)
}
