package repository

import (
	"database/sql"
	"fmt"

	"github.com/DeepanshuMishraa/gotodo/internals/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(email string, hashedPassword string) (*models.User, error) {
	query :=
		`
	INSERT INTO users(email,password,created_at,updated_at) VALUES($1,$2,NOW(),NOW())

	RETURNING id, email, created_at, updated_at
	`

	user := &models.User{}

	err := r.db.QueryRow(query, email, hashedPassword).Scan(
		&user.Id,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query :=
		`
	SELECT id, email, password, created_at, updated_at FROM users WHERE email=$1
	`

	user := &models.User{}

	err := r.db.QueryRow(query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil

}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}