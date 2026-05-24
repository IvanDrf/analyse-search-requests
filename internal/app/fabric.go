package app

import (
	"sync"

	"github.com/IvanDrf/analyse-search-requests/internal/config"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/ports/messaging"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/ports/repo"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/ports/service"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/rules"
	"github.com/IvanDrf/analyse-search-requests/internal/infrastructure/messaging/rabbitmq"
	"github.com/IvanDrf/analyse-search-requests/internal/infrastructure/persistence/redis"
	"github.com/IvanDrf/analyse-search-requests/internal/interfaces/http"

	s "github.com/IvanDrf/analyse-search-requests/internal/infrastructure/service"
)

type fabric struct {
	config *config.Config
}

func (f *fabric) NewApp() *App {
	searchRepo := f.newRepo()
	searchService := f.newSearchService(searchRepo, f.newMessageValidator())

	server := f.newSearchServer(searchService)
	consumer := f.newConsumer(searchService)

	return &App{
		server:   server,
		consumer: consumer,
	}
}

func (f *fabric) newSearchServer(searchService service.SearchService) *http.SearchServer {
	handlers := http.NewHandlers(searchService, f.config.App.RequestTime)

	return http.NewSearchServer(f.config.App.Host, f.config.App.Port, handlers)
}

func (f *fabric) newConsumer(serachService service.SearchService) messaging.MessageConsumer {
	conn, ch, queue := rabbitmq.Connect(&f.config.Broker)

	return rabbitmq.NewSearchConsumer(conn, ch, queue, serachService, f.config.Broker.Workers)
}

func (f *fabric) newSearchService(messageRepo repo.MessageRepo, validator *rules.MessageValidator) service.SearchService {
	return s.NewSearchService(f.config.App.SearchInterval, messageRepo, validator)
}

func (f *fabric) newMessageValidator() *rules.MessageValidator {
	return rules.NewMessageValidator()
}

func (f *fabric) newRepo() repo.MessageRepo {
	conn := redis.Connect(&f.config.Database)
	mx := new(sync.Mutex)

	return redis.NewRedisRepo(
		conn,
		f.config.Database.BadWordKey,
		f.config.App.SearchDuration,
		f.config.Database.DuplicateTime,
		mx,
	)
}
