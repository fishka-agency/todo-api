package service

import (
	"context"

	"github.com/pers0na2dev/todo-api/internal/models"
	"go.uber.org/zap"
)

// TaskRepository интерфейс, который содержит методы для работы с задачами
// Интерфейс позволяет использовать разные реализации репозитория для задач
// Т.е. можно будет легко заменить один репозиторий на другой, например, перейти с Postgres на MongoDB
type TaskRepository interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTasks(ctx context.Context) ([]*models.Task, error)
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id int) error
}

// TaskService структура, которая содержит методы для работы с задачами
type TaskService struct {
	repo   TaskRepository
	logger *zap.Logger
}

// NewTaskService функция, которая создает новый экземпляр TaskService
// @param repo *postgres.TaskRepository - репозиторий для задач
// @return *TaskService - новый экземпляр TaskService
func NewTaskService(repo TaskRepository, logger *zap.Logger) *TaskService {
	return &TaskService{
		repo:   repo,
		logger: logger,
	}
}

// CreateTask функция, которая создает новую задачу
// @param ctx context.Context - контекст выполнения
// @param task *models.Task - задача
// @return error - ошибка
func (s *TaskService) CreateTask(ctx context.Context, title string) error {
	s.logger.Info("creating new task", zap.String("title", title))

	task := &models.Task{Title: title}
	err := s.repo.CreateTask(ctx, task)
	if err != nil {
		s.logger.Error("failed to create task", zap.Error(err))
		return err
	}

	s.logger.Info("task created successfully")
	return nil
}

// GetTasks функция, которая возвращает все задачи
// @param ctx context.Context - контекст выполнения
// @return []*models.Task - список задач
// @return error - ошибка
func (s *TaskService) GetTasks(ctx context.Context) ([]*models.Task, error) {
	s.logger.Info("getting all tasks")

	tasks, err := s.repo.GetTasks(ctx)
	if err != nil {
		s.logger.Error("failed to get tasks", zap.Error(err))
		return nil, err
	}

	s.logger.Info("tasks retrieved successfully", zap.Any("tasks", tasks))
	return tasks, nil
}

// GetTask функция, которая возвращает задачу по id
// @param ctx context.Context - контекст выполнения
// @param id int - id задачи
// @return *models.Task - задача
// @return error - ошибка
func (s *TaskService) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	s.logger.Info("getting task by id", zap.Int("id", id))

	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get task by id", zap.Error(err))
		return nil, err
	}

	s.logger.Info("task retrieved successfully", zap.Any("task", task))
	return task, nil
}

// OpenCloseTask функция, которая открывает или закрывает задачу
// @param ctx context.Context - контекст выполнения
// @param task *models.Task - задача
// @return error - ошибка
func (s *TaskService) OpenCloseTask(ctx context.Context, id int) error {
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get task by id", zap.Error(err))
		return err
	}

	task.Completed = !task.Completed

	err = s.repo.UpdateTask(ctx, task)
	if err != nil {
		s.logger.Error("failed to update task", zap.Error(err))
		return err
	}

	s.logger.Info("task updated successfully")
	return nil
}

// RemoveTask функция, которая удаляет задачу
// @param ctx context.Context - контекст выполнения
// @param id int - id задачи
// @return error - ошибка
func (s *TaskService) RemoveTask(ctx context.Context, id int) error {
	s.logger.Info("removing task", zap.Int("id", id))

	err := s.repo.DeleteTask(ctx, id)
	if err != nil {
		s.logger.Error("failed to remove task", zap.Error(err))
		return err
	}

	s.logger.Info("task removed successfully")
	return nil
}
