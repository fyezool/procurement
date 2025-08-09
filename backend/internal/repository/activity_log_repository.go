package repository

import (
	"database/sql"
	"procurement-system/internal/models"
)

// ActivityLogRepository defines the interface for activity log database operations.
type ActivityLogRepository interface {
	Log(activity *models.ActivityLog) error
	GetAll() ([]models.ActivityLog, error)
}

type postgresActivityLogRepository struct {
	db *sql.DB
}

// NewPostgresActivityLogRepository creates a new instance of ActivityLogRepository.
func NewPostgresActivityLogRepository(db *sql.DB) ActivityLogRepository {
	return &postgresActivityLogRepository{db: db}
}

// Log creates a new activity log entry in the database.
func (r *postgresActivityLogRepository) Log(activity *models.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (user_id, action, target_type, target_id, status, details)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(
		query,
		activity.UserID,
		activity.Action,
		activity.TargetType,
		activity.TargetID,
		activity.Status,
		activity.Details,
	)
	return err
}

// GetAll retrieves all activity logs from the database, ordered by creation date.
func (r *postgresActivityLogRepository) GetAll() ([]models.ActivityLog, error) {
	query := `
		SELECT id, user_id, action, target_type, target_id, status, details, created_at
		FROM activity_logs
		ORDER BY created_at DESC
		LIMIT 100 -- Limit to the last 100 activities for performance
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var log models.ActivityLog
		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.TargetType,
			&log.TargetID,
			&log.Status,
			&log.Details,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
