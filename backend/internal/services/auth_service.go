package services

import (
	"errors"
	"os"
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService interface {
	Register(payload models.RegistrationPayload) (*models.User, error)
	Login(payload models.LoginPayload) (string, error)
}

type authService struct {
	userRepo    repository.UserRepository
	logService  ActivityLogService
}

func NewAuthService(userRepo repository.UserRepository, logService ActivityLogService) AuthService {
	return &authService{userRepo: userRepo, logService: logService}
}

func (s *authService) Register(payload models.RegistrationPayload) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:           payload.Name,
		Email:          payload.Email,
		HashedPassword: string(hashedPassword),
		Role:           payload.Role,
	}

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		details := err.Error()
		s.logService.Log(nil, "REGISTER_USER_FAILED", nil, nil, "FAILED", &details)
		return nil, err
	}

	// Log successful registration
	s.logService.Log(&createdUser.ID, "REGISTER_USER_SUCCESS", Ptr("user"), &createdUser.ID, "SUCCESS", nil)

	// Do not expose password hash in the response
	createdUser.HashedPassword = ""
	return createdUser, nil
}

func (s *authService) Login(payload models.LoginPayload) (string, error) {
	user, err := s.userRepo.GetUserByEmail(payload.Email)
	if err != nil {
		var details string
		if errors.Is(err, repository.ErrUserNotFound) {
			details = "User not found for email: " + payload.Email
			s.logService.Log(nil, "LOGIN_FAILED", Ptr("user"), nil, "FAILED", &details)
			return "", ErrInvalidCredentials
		}
		// Generic database error
		details = err.Error()
		s.logService.Log(nil, "LOGIN_FAILED_DB_ERROR", nil, nil, "FAILED", &details)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(payload.Password))
	if err != nil {
		details := "Invalid password for user: " + payload.Email
		s.logService.Log(&user.ID, "LOGIN_FAILED", Ptr("user"), &user.ID, "FAILED", &details)
		return "", ErrInvalidCredentials
	}

	s.logService.Log(&user.ID, "LOGIN_SUCCESS", Ptr("user"), &user.ID, "SUCCESS", nil)
	return s.generateJWT(user)
}

// Ptr is a helper function to get a pointer to a string.
func Ptr(s string) *string {
	return &s
}

func (s *authService) generateJWT(user *models.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET environment variable not set")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
