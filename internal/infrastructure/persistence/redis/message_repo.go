package redis

import (
	"context"
	"fmt"
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

	lock *sync.Mutex
}

func NewRedisRepo(client *redis.Client, badWordKey string, expTime time.Duration, duplicateTime time.Duration, lock *sync.Mutex) *redisRepo {
	return &redisRepo{
		client:        client,
		badWordKey:    badWordKey,
		expTime:       expTime,
		duplicateTime: duplicateTime,
		lock:          lock,
	}
}

func (r *redisRepo) Close() {
	r.client.Close()
}

func (r *redisRepo) SaveSearch(ctx context.Context, message *models.Message) error {
	key := generateMinuteKey(message.Date)
	duplicate := generateDuplicateKey(message)

	res, err := r.client.Exists(ctx, duplicate).Result()
	if res != 0 && err == nil {
		return nil
	}

	if err := r.client.ZIncrBy(ctx, key, 1, message.SearchMessage); err != nil {
		return models.Error{
			Message: "can't save new search message",
			Code:    models.ErrCodeInternal,
		}
	}

	if err := r.client.Expire(ctx, key, r.expTime).Err(); err != nil {
		return models.Error{
			Message: "can't add exp time for search message",
			Code:    models.ErrCodeInternal,
		}
	}

	if err := r.client.Set(ctx, duplicate, true, r.duplicateTime).Err(); err != nil {
		return models.Error{
			Message: "can't set exp time for duplicate key",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *redisRepo) FindMostPopularSearches(ctx context.Context, limit int, start time.Time, interval int) ([]*models.SearchMessage, error) {
	minutes := generateRangeMinutesKey(start, interval)
	temp := generateTempKey(start)

	r.lock.Lock()
	defer r.lock.Unlock()

	if err := r.client.ZUnionStore(ctx, temp, &redis.ZStore{
		Keys:      minutes,
		Aggregate: "SUM",
	}).Err(); err != nil {
		return nil, models.Error{
			Message: "can't create union storage",
			Code:    models.ErrCodeInternal,
		}
	}

	defer func() {
		r.client.Del(context.Background(), temp)
	}()

	size := limit * 2
	raws, err := r.client.ZRevRangeWithScores(ctx, temp, 0, int64(size-1)).Result()
	if err != nil {
		return nil, models.Error{
			Message: "can't make range with scores",
			Code:    models.ErrCodeInternal,
		}
	}

	return r.filterSearches(ctx, raws, limit)
}

func (r *redisRepo) SaveBadWord(ctx context.Context, badWord string) error {
	if err := r.client.SAdd(ctx, r.badWordKey, badWord).Err(); err != nil {
		return models.Error{
			Message: "can't save bad word in redis",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *redisRepo) DeleteBadWord(ctx context.Context, badWord string) error {
	if err := r.client.SRem(ctx, r.badWordKey, badWord).Err(); err != nil {
		return models.Error{
			Message: "can't delete bad word in redis",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *redisRepo) filterSearches(ctx context.Context, raws []redis.Z, limit int) ([]*models.SearchMessage, error) {
	res := make([]*models.SearchMessage, 0, limit)
	for i := range raws {
		searchMessage, ok := raws[i].Member.(string)
		if !ok {
			return nil, models.Error{
				Message: "invalid searches in redis",
				Code:    models.ErrCodeInternal,
			}
		}

		bad, err := r.client.SIsMember(ctx, r.badWordKey, searchMessage).Result()
		if err != nil {
			return nil, models.Error{
				Message: "can't check stop word in redis",
				Code:    models.ErrCodeInternal,
			}
		}

		if bad {
			continue
		}

		res = append(res, &models.SearchMessage{
			SearchMessage: searchMessage,
			Amount:        uint16(raws[i].Score),
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
	return fmt.Sprintf("search:%s", date.Format("202605231530"))
}

func generateDuplicateKey(message *models.Message) string {
	return fmt.Sprintf("duplicate:%s", message.ID)
}
