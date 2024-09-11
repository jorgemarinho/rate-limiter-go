package web

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jorgemarinho/rate-limiter-go/internal/infra/store"
)

type Limiter struct {
	IPMaxRequestsPerSecond    int
	IPBlockDurationSeconds    int
	TokenMaxRequestsPerSecond int
	TokenBlockDurationSeconds int
	IPStore                   store.Store
	TokenStore                store.Store
}

func NewLimiter(ipMaxRequestsPerSecond, ipBlockDurationSeconds, tokenMaxRequestsPerSecond, tokenBlockDurationSeconds int, ipStore, tokenStore Store) *Limiter {
	return &Limiter{
		IPMaxRequestsPerSecond:    ipMaxRequestsPerSecond,
		IPBlockDurationSeconds:    ipBlockDurationSeconds,
		TokenMaxRequestsPerSecond: tokenMaxRequestsPerSecond,
		TokenBlockDurationSeconds: tokenBlockDurationSeconds,
		IPStore:                   ipStore,
		TokenStore:                tokenStore,
	}
}

func (l *Limiter) AllowRequest(ip, token string) bool {
	if token != "" {
		return l.isAllowedByToken(token)
	}
	return l.isAllowedByIP(ip)
}

func (l *Limiter) isAllowedByIP(ip string) bool {
	key := fmt.Sprintf("ratelimit:ip:%s", ip)
	log.Printf("Checking rate limit for IP: %s", ip)
	return l.isAllowed(key, l.IPMaxRequestsPerSecond, l.IPBlockDurationSeconds, l.IPStore)
}

func (l *Limiter) isAllowedByToken(token string) bool {
	key := fmt.Sprintf("ratelimit:token:%s", token)
	log.Printf("Checking rate limit for Token: %s", token)
	return l.isAllowed(key, l.TokenMaxRequestsPerSecond, l.TokenBlockDurationSeconds, l.TokenStore)
}

func (l *Limiter) isAllowed(key string, maxRequestsPerSecond, blockDurationSeconds int, store store.Store) bool {
	countStr, err := store.Get(key)
	if err != nil {
		log.Printf("Error getting value from store: %v", err)
		return false
	}

	count, _ := strconv.Atoi(countStr)
	log.Printf("Key: %s, Current Count: %d, Max Requests: %d", key, count, maxRequestsPerSecond)
	if count >= maxRequestsPerSecond {
		log.Printf("Rate limit exceeded for key: %s", key)
		return false
	}

	expiration := time.Duration(blockDurationSeconds) * time.Second
	if err := store.Incr(key); err != nil {
		log.Printf("Error incrementing key: %s, error: %v", key, err)
		return false
	}
	if err := store.Expire(key, expiration); err != nil {
		log.Printf("Error setting expiration for key: %s, error: %v", key, err)
		return false
	}

	log.Printf("Request allowed for key: %s", key)
	return true
}
