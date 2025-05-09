package app

import (
	"context"
	"log/slog"

	"github.com/4aykovski/tcp-kursach/internal/config"
	"github.com/4aykovski/tcp-kursach/internal/server"
	"golang.org/x/sync/errgroup"
)

type App struct {
	server *server.Server

	cfg *config.C
}

func New(ctx context.Context) *App {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		panic(err)
	}

	return a
}

func (a *App) Run() error {
	var eg errgroup.Group

	eg.Go(func() error {
		slog.Info("starting tcp server, port", slog.Int("port", a.cfg.Port))
		return a.server.Run(context.Background(), a.cfg)
	})

	return eg.Wait()
}

func (a *App) Stop() error {
	return a.server.Stop(context.Background())
}

func (a *App) initDeps(ctx context.Context) error {
	deps := []func(context.Context) error{
		a.initConfig,
		a.initServer,
	}

	for _, dep := range deps {
		if err := dep(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	a.cfg = config.MustLoad()

	return nil
}

func (a *App) initServer(ctx context.Context) error {
	a.server = server.New(ctx)

	return nil
}
