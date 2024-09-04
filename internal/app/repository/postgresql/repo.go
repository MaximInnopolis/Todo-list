package postgresql

import (
	"context"
	"errors"

	"Todo-list/internal/app/models"
	"Todo-list/internal/app/repository/database"
	"github.com/jackc/pgx/v4"
)

var ErrTaskNotFound = errors.New("task not found")

// Repository interface defines methods for interacting with tasks in database
type Repository interface {
	Create(task models.Task) (models.Task, error)
	GetAll() ([]models.Task, error)
	GetById(id int) (models.Task, error)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}

// Repo struct implements Repository interface and interacts with postgresql database using connection pool
type Repo struct {
	db database.Database
}

// New creates new Repo instance, taking database connection pool as parameter
func New(db database.Database) *Repo {
	return &Repo{db: db}
}

// Create inserts new task into database and returns created task with generated ID, created_at, and updated_at fields
func (r *Repo) Create(task models.Task) (models.Task, error) {
	query := `INSERT INTO tasks (title, description, due_date, created_at, updated_at) 
	          VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at`
	ctx := context.Background()

	// Execute query and scan returned ID, created_at, and updated_at into task object
	err := r.db.GetPool().QueryRow(ctx, query, task.Title, task.Description, task.DueDate).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

// GetAll retrieves all tasks from database and returns them as slice of Task objects
func (r *Repo) GetAll() ([]models.Task, error) {
	query := `SELECT id, title, description, due_date, created_at, updated_at FROM tasks`
	var tasks []models.Task
	ctx := context.Background()

	// Execute query and iterate over result rows
	rows, err := r.db.GetPool().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan each row into Task object and append to tasks slice
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Check for any error that occurred during iteration over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetById retrieves task by ID from database. If task not found, returns ErrTaskNotFound
func (r *Repo) GetById(id int) (models.Task, error) {
	query := `SELECT id, title, description, due_date, created_at, updated_at FROM tasks WHERE id = $1`
	var task models.Task
	ctx := context.Background()

	// Execute query and scan result into Task object
	err := r.db.GetPool().QueryRow(ctx, query, id).
		Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		// If no rows returned, return ErrTaskNotFound.
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}
	return task, nil
}

// Update modifies existing task in database by ID, and returns updated task
// If task with given ID not found, returns ErrTaskNotFound
func (r *Repo) Update(id int, task models.Task) (models.Task, error) {
	query := `UPDATE tasks SET title = $1, description = $2, due_date = $3, updated_at = NOW() 
	          WHERE id = $4 RETURNING id, title, description, due_date, created_at, updated_at`
	ctx := context.Background()

	// Execute query and scan result into task object
	err := r.db.GetPool().QueryRow(ctx, query, task.Title, task.Description, task.DueDate, id).
		Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		// If no rows returned, return ErrTaskNotFound
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}
	return task, nil
}

// Delete removes task from database by ID
// If task with given ID not found, returns ErrTaskNotFound
func (r *Repo) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	ctx := context.Background()

	// Execute delete query and check how many rows were affected
	result, err := r.db.GetPool().Exec(ctx, query, id)
	if err != nil {
		return err
	}

	// If no rows affected, return ErrTaskNotFound
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
