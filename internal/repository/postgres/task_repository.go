package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pers0na2dev/todo-api/internal/models"
	"github.com/pers0na2dev/todo-api/pkg/cache"
	"go.uber.org/zap"
)

const (
	taskCacheKeyPrefix = "task:"
	tasksCacheKey      = "tasks:all"
	cacheDuration      = 5 * time.Minute
)

// TaskRepository структура, которая содержит подключение к базе данных
type TaskRepository struct {
	pool   *pgxpool.Pool
	cache  cache.Cache
	logger *zap.Logger
}

// NewTaskRepository функция, которая создает новый экземпляр TaskRepository
// @param pool *pgxpool.Pool - подключение к базе данных
// @param cache cache.Cache - кеш
// @param logger *zap.Logger - логгер
// @return *TaskRepository - новый экземпляр TaskRepository
func NewTaskRepository(pool *pgxpool.Pool, cache cache.Cache, logger *zap.Logger) *TaskRepository {
	return &TaskRepository{
		pool:   pool,
		cache:  cache,
		logger: logger,
	}
}

// CreateTask функция, которая создает новую задачу
// @param ctx context.Context - контекст выполнения
// @param task *models.Task - задача
// @return error - ошибка
func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING id`
	err := r.pool.QueryRow(ctx, query, task.Title, task.Completed).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// Инвалидируем кеш списка всех задач
	if err := r.cache.Delete(ctx, tasksCacheKey); err != nil {
		r.logger.Warn("failed to invalidate tasks cache", zap.Error(err))
	}

	return nil
}

// GetTasks функция, которая возвращает все задачи
// @param ctx context.Context - контекст выполнения
// @return []*models.Task - список задач
// @return error - ошибка
func (r *TaskRepository) GetTasks(ctx context.Context) ([]*models.Task, error) {
	var tasks []*models.Task

	// Пробуем получить из кеша
	err := r.cache.Get(ctx, tasksCacheKey, &tasks)
	if err == nil {
		r.logger.Debug("tasks retrieved from cache")
		return tasks, nil
	}

	// Если в кеше нет, получаем из БД
	query := `SELECT id, title, completed FROM tasks`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &models.Task{}
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	// Сохраняем в кеш
	if err := r.cache.Set(ctx, tasksCacheKey, tasks, cacheDuration); err != nil {
		r.logger.Warn("failed to cache tasks", zap.Error(err))
	}

	return tasks, nil
}

// GetTaskByID функция, которая возвращает задачу по id
// @param ctx context.Context - контекст выполнения
// @param id int - id задачи
// @return *models.Task - задача
// @return error - ошибка
func (r *TaskRepository) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	task := &models.Task{}
	cacheKey := fmt.Sprintf("%s%d", taskCacheKeyPrefix, id)

	// Пробуем получить из кеша
	err := r.cache.Get(ctx, cacheKey, task)
	if err == nil {
		r.logger.Debug("task retrieved from cache", zap.Int("id", id))
		return task, nil
	}

	// Если в кеше нет, получаем из БД
	query := `SELECT id, title, completed FROM tasks WHERE id = $1`
	err = r.pool.QueryRow(ctx, query, id).Scan(&task.ID, &task.Title, &task.Completed)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Сохраняем в кеш
	if err := r.cache.Set(ctx, cacheKey, task, cacheDuration); err != nil {
		r.logger.Warn("failed to cache task", zap.Error(err))
	}

	return task, nil
}

// UpdateTask функция, которая обновляет задачу
// @param ctx context.Context - контекст выполнения
// @param task *models.Task - задача
// @return error - ошибка
func (r *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks SET title = $1, completed = $2 WHERE id = $3`
	result, err := r.pool.Exec(ctx, query, task.Title, task.Completed, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	// Инвалидируем кеши
	cacheKey := fmt.Sprintf("%s%d", taskCacheKeyPrefix, task.ID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("failed to invalidate task cache", zap.Error(err))
	}
	if err := r.cache.Delete(ctx, tasksCacheKey); err != nil {
		r.logger.Warn("failed to invalidate tasks cache", zap.Error(err))
	}

	return nil
}

// DeleteTask функция, которая удаляет задачу
// @param ctx context.Context - контекст выполнения
// @param id int - id задачи
// @return error - ошибка
func (r *TaskRepository) DeleteTask(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	// Инвалидируем кеши
	cacheKey := fmt.Sprintf("%s%d", taskCacheKeyPrefix, id)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("failed to invalidate task cache", zap.Error(err))
	}
	if err := r.cache.Delete(ctx, tasksCacheKey); err != nil {
		r.logger.Warn("failed to invalidate tasks cache", zap.Error(err))
	}

	return nil
}
