package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func getIpClient(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}

func RateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	client, exits := clients[ip]
	if !exits {
		limiter := rate.NewLimiter(5, 15)
		newClient := &Client{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}
	client.LastSeen = time.Now()
	return client.Limiter
}

func CleanupClient() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.LastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func LimiterMiddleware(rateLimiter zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getIpClient(ctx)
		limiter := RateLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, "You sent too many requests,please try again later!")
			return
		}
		ctx.Next()
	}
}
