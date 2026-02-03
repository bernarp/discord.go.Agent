package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func (a *App) WaitGracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit
	a.log.Info("shutdown signal received", zap.String("signal", s.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.api.Shutdown(ctx); err != nil {
		a.log.Error("failed to shutdown api server", zap.Error(err))
	}

	if err := a.client.Disconnect(); err != nil {
		a.log.Error("failed to disconnect client", zap.Error(err))
	}

	a.log.Info("application stopped")
}
