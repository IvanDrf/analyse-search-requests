package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanDrf/analyse-search-requests/internal/app"
	"github.com/IvanDrf/analyse-search-requests/internal/config"
	"github.com/IvanDrf/analyse-search-requests/internal/infrastructure/adapters"
)

func main() {
	configPath := ""
	flag.StringVar(&configPath, "config", "config/config.example.yaml", "path to config file")
	flag.Parse()

	config := config.LoadFromYaml(configPath)
	adapters.InitLogger(config)

	app := app.NewApp(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.Info("Start service", slog.Int("port", config.App.Port))
	app.Run(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGABRT)
	<-stop

	slog.Info("Stop service", slog.Int("port", config.App.Port))
	app.Stop(ctx)
}
