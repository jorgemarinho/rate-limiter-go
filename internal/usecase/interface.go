package usecase

import (
	"context"

	"github.com/jorgemarinho/rate-limiter-go/internal/entity"
)

type RateLimiterInterface interface {
	VerifyKeyBlock(ctx context.Context, key string) (bool, error)
	BlockKeyPerTime(ctx context.Context, key string, duration int, time string) error
	SetLimitForKeyPerTime(ctx context.Context, key string, duration int, time string) (entity.LimitResult, error)
}
