// use : make sql queries
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/DeepanshuMishraa/api-go/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userId string) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO todos (title, completed, user_id) VALUES ($1, $2, $3) RETURNING id, title, completed, created_at, updated_at,user_id`

	todo := &models.Todo{}
	err := pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func GetAllTodos(pool *pgxpool.Pool, userId string) ([]models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM todos WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := pool.Query(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	todos := []models.Todo{}
	for rows.Next() {
		todo := &models.Todo{}
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserID,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, *todo)
	}

	return todos, nil
}

func GetTodoByID(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM todos WHERE id = $1 ORDER BY created_at DESC`

	todo := &models.Todo{}

	err := pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func UpdateTodo(pool *pgxpool.Pool, title string, completed bool, id int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE todos SET title = $1, completed = $2 , updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id , title, completed, created_at, updated_at `

	todo := &models.Todo{}

	err := pool.QueryRow(ctx, query, title, completed, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM todos WHERE id = $1`

	commandTag, err := pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Todo with id %d not found", id)
	}

	return nil
}
