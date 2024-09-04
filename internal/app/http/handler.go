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

// Handler struct wraps service interface, which interacts with business logic
type Handler struct {
	service api.Service
}

// New creates new Handler instance and takes api.Service as parameter
func New(service api.Service) *Handler {
	return &Handler{service: service}
}

// CreateTaskHandler handles HTTP POST request to create new task
// It parses request body, validates data, and calls service to create task
func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}

	// Decode request body into input struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusBadRequest)
		return
	}

	// Parse due date from string to time.Time format
	dueDate, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		http.Error(w, "Неправильный формат даты", http.StatusBadRequest)
		return
	}

	// Create new task object
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     dueDate,
	}

	// Call service to create task
	createdTask, err := h.service.CreateTask(task)
	if err != nil {
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	// Respond with created task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

// GetTaskListHandler handles HTTP GET request to retrieve all tasks
// It calls service to get list of tasks and returns them in response
func (h *Handler) GetTaskListHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetTaskList()
	if err != nil {
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	// Respond with list of tasks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByIdHandler handles HTTP GET request to retrieve task by ID
// It parses task ID from URL, calls service to get task, and returns task
func (h *Handler) GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert ID string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неправильный формат ID", http.StatusBadRequest)
		return
	}

	// Call service to get task by ID
	task, err := h.service.GetTaskById(id)
	if err != nil {
		// Return 404 error if task is not found
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}

		// Return 500 error for any other issue
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	// Respond with found task
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateTaskByIdHandler handles HTTP PUT request to update existing task by ID
// It parses task ID and input data, calls service to update task, and returns updated task
func (h *Handler) UpdateTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert ID string to integer
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

	// Decode request body into input struct
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusBadRequest)
		return
	}

	// Parse due date from string to time.Time format
	dueDate, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		http.Error(w, "Неправильный формат даты", http.StatusBadRequest)
		return
	}

	// Create new task object with updated data
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     dueDate,
	}

	// Call service to update task by ID
	updatedTask, err := h.service.UpdateTaskById(id, task)
	if err != nil {
		// Return 404 error if task is not found
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}

		// Return 500 error for any other issue
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	// Respond with updated task
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

// DeleteTaskByIdHandler handles HTTP DELETE request to delete task by ID
// It parses task ID, calls service to delete task, and returns appropriate status
func (h *Handler) DeleteTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert ID string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неправильный формат ID", http.StatusBadRequest)
		return
	}

	// Call service to delete task by ID
	err = h.service.DeleteTaskById(id)
	if err != nil {
		// Return 404 error if task is not found
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}
		// Return 500 error for any other issue
		http.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	// Respond with 204 status indicating successful deletion
	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes registers HTTP routes for task operations
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", h.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks", h.GetTaskListHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", h.GetTaskByIdHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", h.UpdateTaskByIdHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", h.DeleteTaskByIdHandler).Methods("DELETE")
}
