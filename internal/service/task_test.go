package service

import (
	"context"
	"errors"
	"testing"

	"github.com/pers0na2dev/todo-api/internal/models"
	"github.com/pers0na2dev/todo-api/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// setupTest подготавливает все необходимые зависимости для тестов
func setupTest(t *testing.T) (*TaskService, *mocks.TaskRepository) {
	t.Helper()

	// Создаем тестовый логгер
	logger, _ := zap.NewDevelopment()

	// Создаем мок репозитория
	mockRepo := new(mocks.TaskRepository)

	// Создаем сервис с моком репозитория
	service := &TaskService{
		repo:   mockRepo,
		logger: logger,
	}

	return service, mockRepo
}

// TestCreateTask тестирует создание новой задачи
func TestCreateTask(t *testing.T) {
	// Подготавливаем тестовые данные
	service, mockRepo := setupTest(t)
	ctx := context.Background()
	title := "Тестовая задача"

	t.Run("Успешное создание задачи", func(t *testing.T) {
		// Настраиваем ожидаемое поведение мока
		mockRepo.On("CreateTask", ctx, &models.Task{Title: title}).Return(nil).Once()

		// Вызываем тестируемый метод
		err := service.CreateTask(ctx, title)

		// Проверяем результаты
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при создании задачи", func(t *testing.T) {
		// Настраиваем ожидаемое поведение мока с ошибкой
		expectedError := errors.New("ошибка базы данных")
		mockRepo.On("CreateTask", ctx, &models.Task{Title: title}).Return(expectedError).Once()

		// Вызываем тестируемый метод
		err := service.CreateTask(ctx, title)

		// Проверяем результаты
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetTasks тестирует получение списка задач
func TestGetTasks(t *testing.T) {
	service, mockRepo := setupTest(t)
	ctx := context.Background()

	t.Run("Успешное получение списка задач", func(t *testing.T) {
		// Подготавливаем тестовые данные
		expectedTasks := []*models.Task{
			{ID: 1, Title: "Задача 1"},
			{ID: 2, Title: "Задача 2"},
		}

		// Настраиваем ожидаемое поведение мока
		mockRepo.On("GetTasks", ctx).Return(expectedTasks, nil).Once()

		// Вызываем тестируемый метод
		tasks, err := service.GetTasks(ctx)

		// Проверяем результаты
		assert.NoError(t, err)
		assert.Equal(t, expectedTasks, tasks)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при получении списка задач", func(t *testing.T) {
		// Настраиваем ожидаемое поведение мока с ошибкой
		expectedError := errors.New("ошибка базы данных")
		mockRepo.On("GetTasks", ctx).Return([]*models.Task(nil), expectedError).Once()

		// Вызываем тестируемый метод
		tasks, err := service.GetTasks(ctx)

		// Проверяем результаты
		assert.Error(t, err)
		assert.Nil(t, tasks)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetTaskByID тестирует получение задачи по ID
func TestGetTaskByID(t *testing.T) {
	service, mockRepo := setupTest(t)
	ctx := context.Background()
	taskID := 1

	t.Run("Успешное получение задачи", func(t *testing.T) {
		// Подготавливаем тестовые данные
		expectedTask := &models.Task{ID: taskID, Title: "Тестовая задача"}

		// Настраиваем ожидаемое поведение мока
		mockRepo.On("GetTaskByID", ctx, taskID).Return(expectedTask, nil).Once()

		// Вызываем тестируемый метод
		task, err := service.GetTaskByID(ctx, taskID)

		// Проверяем результаты
		assert.NoError(t, err)
		assert.Equal(t, expectedTask, task)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при получении задачи", func(t *testing.T) {
		// Настраиваем ожидаемое поведение мока с ошибкой
		expectedError := errors.New("задача не найдена")
		mockRepo.On("GetTaskByID", ctx, taskID).Return(nil, expectedError).Once()

		// Вызываем тестируемый метод
		task, err := service.GetTaskByID(ctx, taskID)

		// Проверяем результаты
		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
