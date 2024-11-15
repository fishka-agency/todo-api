package mocks

import (
	"context"

	"github.com/pers0na2dev/todo-api/internal/models"
	"github.com/stretchr/testify/mock"
)

// TaskRepository это автоматически сгенерированный мок для интерфейса TaskRepository
type TaskRepository struct {
	mock.Mock
}

// CreateTask мок для метода CreateTask
func (m *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

// GetTasks мок для метода GetTasks
func (m *TaskRepository) GetTasks(ctx context.Context) ([]*models.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Task), args.Error(1)
}

// GetTaskByID мок для метода GetTaskByID
func (m *TaskRepository) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

// UpdateTask мок для метода UpdateTask
func (m *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

// DeleteTask мок для метода DeleteTask
func (m *TaskRepository) DeleteTask(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
