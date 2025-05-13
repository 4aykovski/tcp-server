package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/4aykovski/tcp-kursach/internal/app"
)

// main запускает приложение
func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	// создаем контекст, с помощью которого потом будет выполняться плавное завершение программы
	ctx := context.Background()

	// создаем экземпляр приложения
	a := app.New(ctx)

	// запускаем приложение в отдельном потоке, не блокируя основной
	go func() {
		if err := a.Run(); err != nil {
			panic(err)
		}
	}()

	// ждем сигнал системы для завершения программы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := a.Stop(); err != nil {
		panic(err)
	}
}
