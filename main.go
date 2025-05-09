package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/4aykovski/tcp-kursach/internal/app"
)

func main() {

	ctx := context.Background()

	a := app.New(ctx)

	go func() {
		if err := a.Run(); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := a.Stop(); err != nil {
		panic(err)
	}
}
