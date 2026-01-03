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

type LimiterConfig struct {
	Rps   rate.Limit    // Requests per second
	Burst int           // Max burst size
	TTL   time.Duration // Time to live for client data
}

func RateLimiter(ip string, cfg LimiterConfig) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	client, exits := clients[ip]
	if !exits {
		limiter := rate.NewLimiter(cfg.Rps, cfg.Burst)
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

func LimiterMiddleware(rateLimiter *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getIpClient(ctx)
		config := LimiterConfig{Rps: 10, Burst: 10, TTL: 10 * time.Second}
		limiter := RateLimiter(ip, config)

		if !limiter.Allow() {
			if shouldLogRateLimit(ip) {
				rateLimiter.Warn().
					Str("method", ctx.Request.Method).
					Str("path", ctx.Request.URL.Path).
					Str("query", ctx.Request.URL.RawQuery).
					Str("client_ip", ctx.ClientIP()).
					Str("user_agent", ctx.Request.UserAgent()).
					Str("refer", ctx.Request.Referer()).
					Str("protocol", ctx.Request.Proto).
					Str("host", ctx.Request.Host).
					Str("remote_addr", ctx.Request.RemoteAddr).
					Str("request_uri", ctx.Request.RequestURI).
					Int64("content_length", ctx.Request.ContentLength).
					Interface("header", ctx.Request.Header).
					Msg("rate limiter exceeded")
			}
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, "You sent too many requests,please try again later!")
			return
		}
		ctx.Next()
	}
}

var rateLimitLogCache = sync.Map{}

const rateLimitLogTTL = 10 * time.Second

func shouldLogRateLimit(ip string) bool {
	now := time.Now()
	if val, ok := rateLimitLogCache.Load(ip); ok {
		if lastLogTime, ok := val.(time.Time); ok && now.Sub(lastLogTime) < rateLimitLogTTL {
			return false
		}
	}
	rateLimitLogCache.Store(ip, now)
	return true
}
