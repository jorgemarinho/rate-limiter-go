package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jorgemarinho/rate-limiter-go/internal/infra/config"
	"github.com/jorgemarinho/rate-limiter-go/internal/infra/web"
	"github.com/jorgemarinho/rate-limiter-go/internal/middleware"
	"github.com/jorgemarinho/rate-limiter-go/internal/usecase"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
}

func main() {
	envs, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	redisRepository := web.NewRedisInteractor(rdb)
	rateLimiterUseCase := usecase.NewRateLimiterUseCase(redisRepository)

	r := gin.Default()
	r.Use(middleware.RateLimiter(rateLimiterUseCase, envs))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
