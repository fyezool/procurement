package repository

import (
	"database/sql"
	"errors"
	"procurement-system/internal/models"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) CreateUser(user *models.User) error {
	// Check if email already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", user.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailExists
	}

	query := `
		INSERT INTO users (name, email, hashed_password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = r.db.QueryRow(query, user.Name, user.Email, user.HashedPassword, user.Role).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresUserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, hashed_password, role
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
