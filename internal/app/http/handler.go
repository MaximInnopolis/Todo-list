package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"Todo-list/internal/app/api"
	"Todo-list/internal/app/models"
	"Todo-list/internal/app/repository/postgresql"
	"github.com/gorilla/mux"
)

type Handler struct {
	service api.Service
}

func New(service api.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusBadRequest)
		return
	}

	dueDate, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		http.Error(w, "Неправильный формат даты", http.StatusBadRequest)
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     dueDate,
	}

	createdTask, err := h.service.CreateTask(task)
	if err != nil {
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h *Handler) GetTaskListHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetTaskList()
	if err != nil {
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неправильный формат ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTaskById(id)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}

		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) UpdateTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неправильный формат ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}

	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusBadRequest)
		return
	}

	dueDate, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		http.Error(w, "Неправильный формат даты", http.StatusBadRequest)
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     dueDate,
	}

	updatedTask, err := h.service.UpdateTaskById(id, task)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}

		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *Handler) DeleteTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неправильный формат ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTaskById(id)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", h.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks", h.GetTaskListHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", h.GetTaskByIdHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", h.UpdateTaskByIdHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", h.DeleteTaskByIdHandler).Methods("DELETE")
}
