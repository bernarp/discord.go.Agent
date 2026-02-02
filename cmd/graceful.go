package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func (a *App) WaitGracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit
	a.log.Info("shutdown signal received", zap.String("signal", s.String()))

	if err := a.client.Disconnect(); err != nil {
		a.log.Error("failed to disconnect client", zap.Error(err))
	}

	a.log.Info("application stopped")
}
