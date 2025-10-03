package repository

import (
	"database/sql"
	"fmt"
	"github.com/DeepanshuMishraa/gotodo/internals/models"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{
		db: db,
	}
}

func (r *TodoRepository) CreateTodo(userId int, title string, description string) (*models.Todo, error) {
	query := `
		INSERT INTO todos (user_id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, false, NOW(), NOW())
		RETURNING id, user_id, title, description, completed, created_at, updated_at
	`

	todo := &models.Todo{}

	err := r.db.QueryRow(query, userId, title, description).Scan(&todo.Id, &todo.UserId, &todo.Title, &todo.Description, &todo.IsCompleted, &todo.CreatedAt)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *TodoRepository) GetByID(id int) (*models.Todo, error) {
	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	todo := &models.Todo{}
	err := r.db.QueryRow(query, id).Scan(
		&todo.Id,
		&todo.UserId,
		&todo.Title,
		&todo.Description,
		&todo.IsCompleted,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("todo not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return todo, nil
}

func (r *TodoRepository) GetAllByUserID(userID int) ([]*models.Todo, error) {
	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		todo := &models.Todo{}
		err := rows.Scan(
			&todo.Id,
			&todo.UserId,
			&todo.Title,
			&todo.Description,
			&todo.IsCompleted,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	return todos, nil
}


func (r *TodoRepository) Delete(id int) error {
	query := `DELETE FROM todos WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}