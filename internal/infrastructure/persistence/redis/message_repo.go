package redis

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	client     *redis.Client
	badWordKey string

	expTime       time.Duration
	duplicateTime time.Duration

	mx *sync.Mutex
}

func NewRedisRepo(client *redis.Client, badWordKey string, expTime time.Duration, duplicateTime time.Duration, mx *sync.Mutex) *redisRepo {
	return &redisRepo{
		client:        client,
		badWordKey:    badWordKey,
		expTime:       expTime,
		duplicateTime: duplicateTime,
		mx:            mx,
	}
}

func (r *redisRepo) Close() {
	r.client.Close()
	slog.Info("RedisRepo:Close", slog.String("status", "successfully closed redis repo"))
}

func (r *redisRepo) SaveSearch(ctx context.Context, message *models.Message) error {
	slog.Info("RedisRepo:SaveSearch", slog.String("search", message.SearchMessage))

	key := generateMinuteKey(message.Date)
	duplicate := generateDuplicateKey(message)

	ok, err := r.client.SetNX(ctx, duplicate, 1, r.duplicateTime).Result()
	if err != nil {
		slog.Error("RedisRepo:SaveSearch", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't set duplicate key",
			Code:    models.ErrCodeInternal,
		}
	}

	if !ok {
		slog.Info("RedisRepo:SaveSearch", slog.String("status", "search already saved"), slog.String("search", message.SearchMessage))
		return nil
	}

	if err := r.client.ZIncrBy(ctx, key, 1, message.SearchMessage).Err(); err != nil {
		r.client.Del(ctx, duplicate)

		slog.Error("RedisRepo:SaveSearch", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't save new search message",
			Code:    models.ErrCodeInternal,
		}
	}

	if err := r.client.Expire(ctx, key, r.expTime).Err(); err != nil {
		slog.Error("RedisRepo:SaveSearch", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't add exp time for search message",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("RedisRepo:SaveSearch", slog.String("search", message.SearchMessage), slog.String("status", "successfully saved search"))
	return nil
}

func (r *redisRepo) FindMostPopularSearches(ctx context.Context, limit int, start time.Time, interval int) ([]*models.SearchMessage, error) {
	slog.Info("RedisRepo:FindMostPopularSearches", slog.Int("limit", limit))

	minutes := generateRangeMinutesKey(start, interval)
	temp := generateTempKey(start)

	r.mx.Lock()
	defer r.mx.Unlock()

	if err := r.client.ZUnionStore(ctx, temp, &redis.ZStore{
		Keys:      minutes,
		Aggregate: "SUM",
	}).Err(); err != nil {
		slog.Error("RedisRepo:FindMostPopularSearches", slog.String("error", err.Error()))
		return nil, models.Error{
			Message: "can't create union storage",
			Code:    models.ErrCodeInternal,
		}
	}

	defer func() {
		r.client.Del(context.Background(), temp)
	}()

	size := limit * 3
	raws, err := r.client.ZRevRangeWithScores(ctx, temp, 0, int64(size-1)).Result()
	if err != nil {
		slog.Error("RedisRepo:FindMostPopularSearches", slog.String("error", err.Error()))
		return nil, models.Error{
			Message: "can't make range with scores",
			Code:    models.ErrCodeInternal,
		}
	}

	return r.filterSearches(ctx, raws, limit)
}

func (r *redisRepo) SaveBadWord(ctx context.Context, badWord string) error {
	if err := r.client.SAdd(ctx, r.badWordKey, badWord).Err(); err != nil {
		slog.Error("RedisRepo:SaveBadWord", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't save bad word in redis",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("RedisRepo:SaveBadWord", slog.String("status", "successfully saved bad word"))
	return nil
}

func (r *redisRepo) DeleteBadWord(ctx context.Context, badWord string) error {
	if err := r.client.SRem(ctx, r.badWordKey, badWord).Err(); err != nil {
		slog.Error("RedisRepo:DeleteBadWord", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't delete bad word in redis",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("RedisRepo:SaveBadWord", slog.String("status", "successfully deleted bad word"))
	return nil
}

func (r *redisRepo) filterSearches(ctx context.Context, raws []redis.Z, limit int) ([]*models.SearchMessage, error) {
	res := make([]*models.SearchMessage, 0, limit)
	for i := range raws {
		searchMessage, ok := raws[i].Member.(string)
		if !ok {
			slog.Error("RedisRepo:filterSearches", slog.String("error", "invalid searches in redis"))
			return nil, models.Error{
				Message: "invalid searches in redis",
				Code:    models.ErrCodeInternal,
			}
		}

		skip, err := r.client.SIsMember(ctx, r.badWordKey, searchMessage).Result()
		if err != nil {
			slog.Error("RedisRepo:filterSearches", slog.String("error", err.Error()))
			return nil, models.Error{
				Message: "can't check stop word in redis",
				Code:    models.ErrCodeInternal,
			}
		}

		for word := range strings.SplitSeq(searchMessage, " ") {
			skip, err = r.client.SIsMember(ctx, r.badWordKey, word).Result()
			if err != nil {
				slog.Error("RedisRepo:filterSearches", slog.String("error", err.Error()))
				return nil, models.Error{
					Message: "can't check stop word in redis",
					Code:    models.ErrCodeInternal,
				}
			}

			if skip {
				break
			}
		}

		if skip {
			continue
		}

		res = append(res, &models.SearchMessage{
			SearchMessage: searchMessage,
			Amount:        uint64(raws[i].Score),
		})

		if len(res) == limit {
			return res, nil
		}
	}

	return res, nil
}

func generateRangeMinutesKey(start time.Time, interval int) []string {
	minutes := make([]string, 0, interval)
	for i := range interval {
		date := start.Add(-time.Minute * time.Duration(i))
		minutes = append(minutes, generateMinuteKey(date))
	}

	return minutes
}

func generateTempKey(start time.Time) string {
	return fmt.Sprintf("temp-key:%d", start.Unix())

}

func generateMinuteKey(date time.Time) string {
	return fmt.Sprintf("search:%s", date.Format("200601021504"))
}

func generateDuplicateKey(message *models.Message) string {
	return fmt.Sprintf("duplicate:%s", message.ID)
}
