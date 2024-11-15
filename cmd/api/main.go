package main

import (
	"github.com/pers0na2dev/todo-api/internal/api"
	"github.com/pers0na2dev/todo-api/internal/api/handlers"
	"github.com/pers0na2dev/todo-api/internal/config"
	"github.com/pers0na2dev/todo-api/internal/repository"
	"github.com/pers0na2dev/todo-api/internal/repository/postgres"
	"github.com/pers0na2dev/todo-api/internal/service"
	"github.com/pers0na2dev/todo-api/pkg/cache"
	"github.com/pers0na2dev/todo-api/pkg/logger"

	"go.uber.org/fx"
)

func main() {
	// fx.New - создает новое приложение с Dependency Injection
	fx.New(
		createApp(), // создание приложения
	).Run()
}

// createApp функция, которая создает приложение
// @return fx.Option - опции для приложения
func createApp() fx.Option {
	return fx.Options(
		// fx.Provide - предоставляет зависимости которые будут использоваться в приложении
		fx.Provide(
			config.LoadConfig,                // загрузка конфигурации
			logger.NewLogger,                 // создание логгера
			repository.NewPostgresConnection, // подключение к базе данных
			fx.Annotate(
				cache.NewRedisCache,     // создание кеша
				fx.As(new(cache.Cache)), // указываем что кеш реализует интерфейс Cache
			),
			fx.Annotate(
				postgres.NewTaskRepository,         // создание репозитория для задач
				fx.As(new(service.TaskRepository)), // указываем что репозиторий для задач реализует интерфейс TaskRepository
			),
			fx.Annotate(
				service.NewTaskService,           // создание сервиса для задач
				fx.As(new(handlers.TaskService)), // указываем что сервис для задач реализует интерфейс TaskService
			),
			api.NewServer, // создание HTTP сервера
		),
		// fx.Invoke - вызывает функции которые будут выполняться при запуске приложения
		fx.Invoke(
			handlers.NewTaskHandler, // создание обработчика для задач
		),
	)
}
