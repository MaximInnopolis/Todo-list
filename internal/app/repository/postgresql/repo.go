package postgresql

import (
	"context"
	"errors"

	"Todo-list/internal/app/models"
	"Todo-list/internal/app/repository/database"
	"github.com/jackc/pgx/v4"
)

var ErrTaskNotFound = errors.New("task not found")

type Repository interface {
	Create(task models.Task) (models.Task, error)
	GetAll() ([]models.Task, error)
	GetById(id int) (models.Task, error)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}

type Repo struct {
	db database.Database
}

func New(db database.Database) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(task models.Task) (models.Task, error) {
	query := `INSERT INTO tasks (title, description, due_date, created_at, updated_at) 
	          VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at`
	ctx := context.Background()

	err := r.db.GetPool().QueryRow(ctx, query, task.Title, task.Description, task.DueDate).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil

}

func (r *Repo) GetAll() ([]models.Task, error) {
	query := `SELECT id, title, description, due_date, created_at, updated_at FROM tasks`
	var tasks []models.Task
	ctx := context.Background()

	rows, err := r.db.GetPool().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repo) GetById(id int) (models.Task, error) {
	query := `SELECT id, title, description, due_date, created_at, updated_at FROM tasks WHERE id = $1`
	var task models.Task
	ctx := context.Background()

	err := r.db.GetPool().QueryRow(ctx, query, id).
		Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}
	return task, nil
}

func (r *Repo) Update(id int, task models.Task) (models.Task, error) {
	query := `UPDATE tasks SET title = $1, description = $2, due_date = $3, updated_at = NOW() 
	          WHERE id = $4 RETURNING id, title, description, due_date, created_at, updated_at`
	ctx := context.Background()

	err := r.db.GetPool().QueryRow(ctx, query, task.Title, task.Description, task.DueDate, id).
		Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}
	return task, nil
}

func (r *Repo) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	ctx := context.Background()

	result, err := r.db.GetPool().Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
