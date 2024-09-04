package api

import (
	"Todo-list/internal/app/models"
	"Todo-list/internal/app/repository/postgresql"
)

type Service interface {
	CreateTask(task models.Task) (models.Task, error)
	GetTaskList() ([]models.Task, error)
	GetTaskById(id int) (models.Task, error)
	UpdateTaskById(id int, task models.Task) (models.Task, error)
	DeleteTaskById(id int) error
}

type TaskService struct {
	repo postgresql.Repository
}

func New(repo postgresql.Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (t *TaskService) CreateTask(task models.Task) (models.Task, error) {
	return t.repo.Create(task)
}

func (t *TaskService) GetTaskList() ([]models.Task, error) {
	return t.repo.GetAll()
}

func (t *TaskService) GetTaskById(id int) (models.Task, error) {
	return t.repo.GetById(id)
}

func (t *TaskService) UpdateTaskById(id int, task models.Task) (models.Task, error) {
	return t.repo.Update(id, task)
}

func (t *TaskService) DeleteTaskById(id int) error {
	return t.repo.Delete(id)
}
