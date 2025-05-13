package config

import "github.com/ilyakaznacheev/cleanenv"

// C структура, описывающая конфиг приложения
type C struct {
	TCP
}

// TCP структура, описывающая конфиг tcp-сервера
type TCP struct {
	Port int    `env:"TCP_SERVER_PORT" env-default:"8000"`
	Host string `env:"TCP_SERVER_HOST" env-default:"127.0.0.1"`
}

// MustLoad инициализирует конфиг приложения
func MustLoad() *C {
	var cfg C

	// загружаем конфиг из переменных окружения
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
