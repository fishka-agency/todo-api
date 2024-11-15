package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pers0na2dev/todo-api/internal/models"
)

// TaskService интерфейс, который определяет методы для работы с задачами
// интерфейс для сервиса задач позволяет использовать разные реализации сервиса задач с одинаковым интерфейсом
type TaskService interface {
	CreateTask(ctx context.Context, title string) error
	GetTasks(ctx context.Context) ([]*models.Task, error)
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)
	OpenCloseTask(ctx context.Context, id int) error
	RemoveTask(ctx context.Context, id int) error
}

type TaskHandler struct {
	taskService TaskService
}

func NewTaskHandler(taskService TaskService, mux *http.ServeMux) *TaskHandler {
	handler := &TaskHandler{taskService: taskService}

	mux.HandleFunc("GET /v1/tasks", handler.GetTasks)
	mux.HandleFunc("GET /v1/tasks/{id}", handler.GetTaskByID)
	mux.HandleFunc("POST /v1/tasks", handler.CreateTask)
	mux.HandleFunc("PUT /v1/tasks/{id}", handler.ChangeTaskStatus)
	mux.HandleFunc("DELETE /v1/tasks/{id}", handler.RemoveTask)

	return &TaskHandler{taskService: taskService}
}

// GetTasks функция, которая возвращает все задачи
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.GetTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Установка заголовка Content-Type для ответа
	w.Header().Set("Content-Type", "application/json")
	// Кодирование задач в JSON и отправка ответа
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByID функция, которая возвращает задачу по id
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// Получение id задачи из URL
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение задачи по id
	task, err := h.taskService.GetTaskByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Установка заголовка Content-Type для ответа
	w.Header().Set("Content-Type", "application/json")
	// Кодирование задачи в JSON и отправка ответа
	json.NewEncoder(w).Encode(task)
}

// CreateTask функция, которая создает новую задачу
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	// Декодирование тела запроса в структуру task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создание задачи в базе данных
	err = h.taskService.CreateTask(r.Context(), task.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа с кодом 201 Created
	w.WriteHeader(http.StatusCreated)
}

// ChangeTaskStatus функция, которая изменяет статус задачи
func (h *TaskHandler) ChangeTaskStatus(w http.ResponseWriter, r *http.Request) {
	// Получение id задачи из URL
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Изменение статуса задачи
	err = h.taskService.OpenCloseTask(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа с кодом 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// RemoveTask функция, которая удаляет задачу
func (h *TaskHandler) RemoveTask(w http.ResponseWriter, r *http.Request) {
	// Получение id задачи из URL
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Удаление задачи из базы данных
	err = h.taskService.RemoveTask(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа с кодом 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
