package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/haneumLee/legacychain/backend/config"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(cfg *config.Config, redisClient *redis.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		ctx := context.Background()

		// Increment counter
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Rate limit error",
			})
		}

		// Set expiry on first request
		if count == 1 {
			redisClient.Expire(ctx, key, cfg.RateLimit.Window)
		}

		// Check limit
		if count > int64(cfg.RateLimit.Max) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		}

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.RateLimit.Max))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", cfg.RateLimit.Max-int(count)))

		ttl, _ := redisClient.TTL(ctx, key).Result()
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))

		return c.Next()
	}
}
