package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanDrf/analyse-search-requests/internal/app"
	"github.com/IvanDrf/analyse-search-requests/internal/config"
)

func main() {
	configPath := ""
	flag.StringVar(&configPath, "config", "config/config.example.yaml", "path to config file")
	flag.Parse()

	config := config.LoadFromYaml(configPath)
	app := app.NewApp(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.Run(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGABRT)
	<-stop

	app.Stop()
}
