package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	// HTTP_PORT - порт для запуска HTTP сервера
	HTTPPort string `mapstructure:"HTTP_PORT"`
	// POSTGRES_DSN - строка подключения к базе данных
	PostgresDSN string `mapstructure:"POSTGRES_DSN"`
	// REDIS_DSN - строка подключения к Redis
	RedisDSN string `mapstructure:"REDIS_DSN"`
}

func LoadConfig() (*Config, error) {
	// Загрузка переменных окружения из файла .env
	viper.SetConfigFile(".env")

	// Автоматически использовать переменные окружения, если они доступны
	viper.AutomaticEnv()

	// Чтение конфигурационного файла
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Преобразование конфигурации в структуру Config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Возврат загруженной конфигурации
	return &cfg, nil
}
