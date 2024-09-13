package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jorgemarinho/rate-limiter-go/internal/infra/config"
	"github.com/jorgemarinho/rate-limiter-go/internal/usecase"
)

func RateLimiter(rateUseCase *usecase.RateLimiterUseCase, envs *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		inputIP := usecase.InputRateLimiter{
			Item:           c.ClientIP(),
			RateLimit:      envs.RateLimitIP,
			TimeTypeLimit:  envs.TimeLimitType,
			BlockLimitTime: envs.BlockLimitTimeDuration,
			TimeTypeBlock:  envs.TimeBlockType,
		}
		outputIP := rateUseCase.Execute(ctx, inputIP)

		if outputIP.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println("Internal Server Error")
			c.Abort()
			return
		}

		inputToken := usecase.InputRateLimiter{
			Item:           c.GetHeader("API_KEY"),
			RateLimit:      envs.RateLimitToken,
			TimeTypeLimit:  envs.TimeLimitType,
			BlockLimitTime: envs.BlockLimitTimeDuration,
			TimeTypeBlock:  envs.TimeBlockType,
		}
		outputToken := rateUseCase.Execute(ctx, inputToken)

		if outputToken.Err != nil {
			if outputToken.Err.Error() == "input empty" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Empty TOKEN"})
				log.Println("Empty TOKEN")
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println("Internal Server Error")
			c.Abort()
			return
		}

		if !outputIP.AllowRequest && !outputToken.AllowRequest {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
			log.Println("Rate limit exceeded for Token and IP:", c.GetHeader("API_KEY"), c.ClientIP())
			c.Abort()
			return
		}

		if !outputToken.AllowRequest {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
			log.Println("Rate limit exceeded for Token:", c.GetHeader("API_KEY"))
			c.Abort()
			return
		}

		if !outputIP.AllowRequest && outputToken.AllowRequest {
			log.Println("Rate limit exceeded for IP but NOT for Token", c.GetHeader("API_KEY"))
		}

		c.Next()
	}
}
