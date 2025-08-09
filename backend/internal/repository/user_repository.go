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
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	UpdatePassword(userID int, newHashedPassword string) error
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) CreateUser(user *models.User) (*models.User, error) {
	// Check if email already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", user.Email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	query := `
		INSERT INTO users (name, email, hashed_password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = r.db.QueryRow(query, user.Name, user.Email, user.HashedPassword, user.Role).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *postgresUserRepository) UpdatePassword(userID int, newHashedPassword string) error {
	query := `UPDATE users SET hashed_password = $1 WHERE id = $2`
	result, err := r.db.Exec(query, newHashedPassword, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *postgresUserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, hashed_password, role
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetAllUsers() ([]models.User, error) {
	query := `
		SELECT id, name, email, role
		FROM users
		ORDER BY name ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *postgresUserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, role = $3
		WHERE id = $4
	`
	result, err := r.db.Exec(query, user.Name, user.Email, user.Role, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *postgresUserRepository) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
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
