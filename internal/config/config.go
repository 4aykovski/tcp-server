package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// C структура, описывающая конфиг приложения
type C struct {
	TCP
}

// TCP структура, описывающая конфиг tcp-сервера
type TCP struct {
	Port int    `env:"TCP_SERVER_PORT" env-required:"true"`
	Host string `env:"TCP_SERVER_HOST" env-required:"true"`
}

// MustLoad инициализирует конфиг приложения
func MustLoad() *C {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	var cfg C

	// загружаем конфиг из переменных окружения
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
