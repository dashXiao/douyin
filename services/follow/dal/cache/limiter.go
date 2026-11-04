package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Limit(ctx context.Context, rate int, interval time.Duration) error {
	key := LimiterKey
	now := time.Now().UnixNano()
	// Remove request timestamps outside the sliding window.
	RedisClient.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(interval)))
	// Count requests in the current window.
	size, err := RedisClient.ZCard(ctx, key).Result()
	if err != nil {
		return err
	}
	// Reject the request when the window has reached its limit.
	if size >= int64(rate) {
		return errors.New("request rate limit exceeded")
	}
	// Record the current request timestamp in the sorted set.
	RedisClient.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: float64(now)})
	return nil
}
