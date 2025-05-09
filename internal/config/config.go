package config

import "github.com/ilyakaznacheev/cleanenv"

type C struct {
	TCP
}

type TCP struct {
	Port int    `env:"TCP_SERVER_PORT" env-default:"8000"`
	Host string `env:"TCP_SERVER_HOST" env-default:"127.0.0.1"`
}

func MustLoad() *C {
	var cfg C

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
