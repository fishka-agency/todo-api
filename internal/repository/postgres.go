package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pers0na2dev/todo-api/internal/config"
	"go.uber.org/zap"
)

// NewPostgresConnection функция, которая создает подключение к базе данных
// @param cfg *config.Config - конфигурация
// @param logger *zap.Logger - логгер
// @return *pgxpool.Pool - подключение к базе данных
// @return error - ошибка
func NewPostgresConnection(cfg *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	logger.Info("connecting to postgres database")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		logger.Error("failed to create connection pool", zap.Error(err))
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("successfully connected to postgres database")
	return pool, nil
}
