package api

import (
	"Todo-list/internal/app/models"
	"Todo-list/internal/app/repository/postgresql"
)

// Service defines interface for task service, which includes methods
// to create, retrieve, update, and delete tasks.
type Service interface {
	CreateTask(task models.Task) (models.Task, error)
	GetTaskList() ([]models.Task, error)
	GetTaskById(id int) (models.Task, error)
	UpdateTaskById(id int, task models.Task) (models.Task, error)
	DeleteTaskById(id int) error
}

// TaskService is implementation of Service interface.
// It interacts with repository to perform CRUD operations on tasks.
type TaskService struct {
	repo postgresql.Repository
}

// New creates new TaskService instance and takes Repository as parameter
func New(repo postgresql.Repository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask creates new task using repository and returns created task
func (t *TaskService) CreateTask(task models.Task) (models.Task, error) {
	return t.repo.Create(task)
}

// GetTaskList retrieves list of all tasks from repository
func (t *TaskService) GetTaskList() ([]models.Task, error) {
	return t.repo.GetAll()
}

// GetTaskById retrieves task by ID using repository
func (t *TaskService) GetTaskById(id int) (models.Task, error) {
	return t.repo.GetById(id)
}

// UpdateTaskById updates an existing task by ID using repository
// and returns updated task
func (t *TaskService) UpdateTaskById(id int, task models.Task) (models.Task, error) {
	return t.repo.Update(id, task)
}

// DeleteTaskById deletes task by ID using repository
func (t *TaskService) DeleteTaskById(id int) error {
	return t.repo.Delete(id)
}
