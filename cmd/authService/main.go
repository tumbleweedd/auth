package main

import (
	"context"
	"github.com/tumbleweedd/svc/auth_service/internal/app"
	"github.com/tumbleweedd/svc/auth_service/internal/config"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.InitConfig()

	log := logger.NewSlogLogger(logger.SlogEnvironment(cfg.Env))

	application := app.New(ctx, cfg, log, cfg.GRPCConfig.Port)

	go func() {
		application.MustRun()
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop(ctx)

	log.Info("Application stopped")
}
