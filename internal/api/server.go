package api

import (
	"context"
	"net/http"

	"github.com/pers0na2dev/todo-api/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewServer создает новый HTTP сервер
func NewServer(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) *http.ServeMux {
	// Создание нового маршрутизатора HTTP, который будет использоваться для обработки запросов
	mux := http.NewServeMux()

	// Создание нового HTTP сервера с заданным адресом и маршрутизатором
	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: mux,
	}

	// Добавление хука жизненного цикла fx для запуска и завершения сервера
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting http server", zap.String("port", cfg.HTTPPort))
			go server.ListenAndServe() // Запуск HTTP сервера в отдельной горутине
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down http server")
			return server.Shutdown(ctx) // Завершение работы HTTP сервера
		},
	})

	return mux
}
