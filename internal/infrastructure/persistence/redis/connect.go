package redis

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/IvanDrf/analyse-search-requests/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect(config *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: config.Addr(),

		Username: config.Username,
		Password: config.Password,

		DB: config.Database,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("can't ping redis database, error=%s", err)
	}

	slog.Info("successfully connected to redis", slog.String("host", config.Host), slog.Int("port", config.Port))
	return client
}
