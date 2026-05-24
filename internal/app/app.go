package app

import (
	"context"

	"github.com/IvanDrf/analyse-search-requests/internal/config"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/ports/messaging"
	"github.com/IvanDrf/analyse-search-requests/internal/interfaces/http"
)

type App struct {
	server   *http.SearchServer
	consumer messaging.MessageConsumer
}

func NewApp(config *config.Config) *App {
	fabric := fabric{config: config}

	return fabric.NewApp()
}

func (a *App) Run(ctx context.Context) {
	go a.server.Start()
	go a.consumer.StartReadingMessages(ctx)
}

func (a *App) Stop() {
	a.server.Stop()
	a.consumer.Close()
}
