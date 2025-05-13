package app

import (
	"context"
	"log/slog"

	"github.com/4aykovski/tcp-kursach/internal/config"
	"github.com/4aykovski/tcp-kursach/internal/server"
	"golang.org/x/sync/errgroup"
)

// App структура, содержащая в себе сам tcp-сервер и конфиг приложения.
type App struct {
	server *server.Server

	cfg *config.C
}

// New создает новый экземпляр приложения.
func New(ctx context.Context) *App {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		panic(err)
	}

	return a
}

// Run запускает приложение
func (a *App) Run() error {
	var eg errgroup.Group

	// запускает tcp-сервер в отдельном потоке, чтобы не блокировать основной
	eg.Go(func() error {
		slog.Info("starting tcp server, port", slog.Int("port", a.cfg.Port))
		return a.server.Run(context.Background(), a.cfg)
	})

	// возвращаем ошибку, если таковая случилась в потоке tcp-сервера
	return eg.Wait()
}

// Stop останавливает приложение
func (a *App) Stop() error {
	// даем команду tcp-серверу, что нужно завершится
	return a.server.Stop(context.Background())
}

// initDeps инициализирует зависимости приложения - сервер и конфиг
func (a *App) initDeps(ctx context.Context) error {
	// список функций, которые необходимо выполнить, чтобы инициализировать зависимости
	deps := []func(context.Context) error{
		a.initConfig,
		a.initServer,
	}

	// проходимся в цикле по функциям, запуская каждую
	for _, dep := range deps {
		if err := dep(ctx); err != nil {
			return err
		}
	}

	return nil
}

// initConfig инициализирует конфиг
func (a *App) initConfig(ctx context.Context) error {
	a.cfg = config.MustLoad()

	return nil
}

// initServer инициализирует tcp-сервер
func (a *App) initServer(ctx context.Context) error {
	a.server = server.New(ctx)

	return nil
}
