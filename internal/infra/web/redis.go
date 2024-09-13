package web

import (
	"context"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/jorgemarinho/rate-limiter-go/internal/entity"
	"github.com/redis/go-redis/v9"
)

type RedisInteractor struct {
	rdb     *redis.Client
	limiter *redis_rate.Limiter
}

func NewRedisInteractor(rdb *redis.Client) *RedisInteractor {
	limiter := redis_rate.NewLimiter(rdb)

	return &RedisInteractor{
		rdb:     rdb,
		limiter: limiter,
	}
}

func (r *RedisInteractor) VerifyKeyBlock(ctx context.Context, key string) (bool, error) {
	_, err := r.rdb.Get(ctx, key+":blocked").Result()
	if err != redis.Nil && err != nil {
		return false, err
	}

	if err == redis.Nil {
		return false, nil
	}

	return true, nil
}

func (r *RedisInteractor) BlockKeyPerTime(ctx context.Context, key string, duration int, t string) error {
	var timeToBlock time.Duration
	switch t {
	case "minute":
		timeToBlock = time.Duration(duration) * time.Minute
	case "second":
		timeToBlock = time.Duration(duration) * time.Second
	case "hour":
		timeToBlock = time.Duration(duration) * time.Hour
	}

	err := r.rdb.Set(ctx, key+":blocked", "true", timeToBlock).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisInteractor) SetLimitForKeyPerTime(ctx context.Context, key string, rate int, time string) (entity.LimitResult, error) {
	var rateForLimiter redis_rate.Limit

	switch time {
	case "minute":
		rateForLimiter = redis_rate.PerMinute(rate)
	case "second":
		rateForLimiter = redis_rate.PerSecond(rate)
	case "hour":
		rateForLimiter = redis_rate.PerHour(rate)
	}

	result, err := r.limiter.Allow(ctx, key, rateForLimiter)
	if err != nil {
		return entity.LimitResult{}, err
	}

	return entity.LimitResult{
		Allowed:   result.Allowed,
		Remaining: result.Remaining,
	}, nil
}
