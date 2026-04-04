package repository

import (
	"context"
	"time"

	"github.com/DeepanshuMishraa/api-go/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id,email, created_at, updated_at`

	createdUsers := &models.User{}

	err := pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&createdUsers.ID,
		&createdUsers.Email,
		&createdUsers.CreatedAt,
		&createdUsers.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return createdUsers, nil

}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`

	users := &models.User{}

	err := pool.QueryRow(ctx, query, email).Scan(
		&users.ID,
		&users.Email,
		&users.Password,
		&users.CreatedAt,
		&users.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return users, nil

}

func GetUserById(pool *pgxpool.Pool, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, email, created_at, updated_at FROM users WHERE id = $1`

	users := &models.User{}

	err := pool.QueryRow(ctx, query, id).Scan(
		&users.ID,
		&users.Email,
		&users.CreatedAt,
		&users.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return users, nil

}
