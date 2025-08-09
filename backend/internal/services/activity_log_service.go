package services

import (
	"log"
	"procurement-system/internal/models"
	"procurement-system/internal/repository"
)

// ActivityLogService defines the interface for activity logging operations.
type ActivityLogService interface {
	Log(userID *int, action string, targetType *string, targetID *int, status string, details *string)
	GetAll() ([]models.ActivityLog, error)
}

type activityLogService struct {
	repo repository.ActivityLogRepository
}

// NewActivityLogService creates a new instance of ActivityLogService.
func NewActivityLogService(repo repository.ActivityLogRepository) ActivityLogService {
	return &activityLogService{repo: repo}
}

// Log logs an activity. It runs in a separate goroutine so it doesn't block the main request flow.
func (s *activityLogService) Log(userID *int, action string, targetType *string, targetID *int, status string, details *string) {
	go func() {
		activity := &models.ActivityLog{
			UserID:     userID,
			Action:     action,
			TargetType: targetType,
			TargetID:   targetID,
			Status:     status,
			Details:    details,
		}
		// We log the error here for observability but don't block the main thread.
		if err := s.repo.Log(activity); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()
}

// GetAll retrieves all activity logs.
func (s *activityLogService) GetAll() ([]models.ActivityLog, error) {
	return s.repo.GetAll()
}
