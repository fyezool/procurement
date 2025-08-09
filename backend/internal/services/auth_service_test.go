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

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)

	payload := models.RegistrationPayload{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "Employee",
	}

	// We can't directly compare the hashed password, so we capture the user argument
	var capturedUser *models.User
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		capturedUser = args.Get(0).(*models.User)
	}).Return(nil)

	user, err := authService.Register(payload)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	mockRepo.AssertExpectations(t)

	// Verify captured user details
	assert.Equal(t, payload.Name, capturedUser.Name)
	assert.Equal(t, payload.Email, capturedUser.Email)
	assert.Equal(t, payload.Role, capturedUser.Role)
	err = bcrypt.CompareHashAndPassword([]byte(capturedUser.HashedPassword), []byte(payload.Password))
	assert.NoError(t, err, "Hashed password should match original password")
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)

	payload := models.RegistrationPayload{
		Email: "exists@example.com",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(repository.ErrEmailExists)

	user, err := authService.Register(payload)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, repository.ErrEmailExists, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUser := &models.User{
		ID:             1,
		Email:          "test@example.com",
		HashedPassword: string(hashedPassword),
		Role:           "Admin",
	}

	mockRepo.On("GetUserByEmail", "test@example.com").Return(mockUser, nil)

	token, err := authService.Login(models.LoginPayload{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)

	// Test case 1: User not found
	mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, repository.ErrUserNotFound)
	_, err := authService.Login(models.LoginPayload{Email: "notfound@example.com", Password: "password"})
	assert.Equal(t, ErrInvalidCredentials, err)

	// Test case 2: Wrong password
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUser := &models.User{Email: "test@example.com", HashedPassword: string(hashedPassword)}

	// We need to set up the mock again for the new test case
	mockRepo.On("GetUserByEmail", "test@example.com").Return(mockUser, nil)
	_, err = authService.Login(models.LoginPayload{Email: "test@example.com", Password: "wrongpassword"})
	assert.Equal(t, ErrInvalidCredentials, err)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)

	expectedErr := errors.New("database error")
	mockRepo.On("GetUserByEmail", "any@example.com").Return(nil, expectedErr)

	_, err := authService.Login(models.LoginPayload{Email: "any@example.com", Password: "password"})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}
