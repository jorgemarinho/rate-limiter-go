package main

import (
	"log"
	"net/http"
	"path/filepath"

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

	// Obter o caminho absoluto do arquivo .env
	envPath, err := filepath.Abs("../../.env")
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	envs, err := config.LoadConfig(envPath)
	if err != nil {
		panic(err)
	}

	redisRepository := web.NewRedisInteractor(rdb)
	rateLimiterUseCase := usecase.NewRateLimiterUseCase(redisRepository)

	r := gin.Default()
	r.Use(middleware.RateLimiter(rateLimiterUseCase, envs))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "welcome to rate limiter"})
	})

	log.Println("Server running on port 8080")

	r.Run(":8080")
}
