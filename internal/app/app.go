package app

import (
	"context"

	"github.com/IvanDrf/analyse-search-requests/internal/config"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/ports/messaging"
	"github.com/IvanDrf/analyse-search-requests/internal/interfaces/http"
)

type App struct {
	searchserver  *http.SearchServer
	metricsServer *http.MetricsServer
	consumer      messaging.MessageConsumer
}

func NewApp(config *config.Config) *App {
	fabric := fabric{config: config}

	return fabric.NewApp()
}

func (a *App) Run(ctx context.Context) {
	go a.metricsServer.Start()
	go a.searchserver.Start()
	go a.consumer.StartReadingMessages(ctx)
}

func (a *App) Stop(ctx context.Context) {
	a.searchserver.Stop(ctx)
	a.metricsServer.Stop(ctx)

	a.consumer.Close()
}
