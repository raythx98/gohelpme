package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/raythx98/gohelpme/tool/logger"

	"golang.org/x/time/rate"
)

type RateConfig struct {
	Rate  float64 `yaml:"rate" json:"rate"`
	Burst int     `yaml:"burst" json:"burst"`
}

type Config struct {
	Default    RateConfig            `yaml:"default" json:"default"`
	Operations map[string]RateConfig `yaml:"operations" json:"operations"`
}

type RateLimiter struct {
	config       Config
	log          logger.ILogger
	limiters     sync.Map
	keyExtractor func(r *http.Request) (ip string, operation string)
	cleanup      *time.Ticker
}

func NewRateLimiter(cfg Config, log logger.ILogger, extractor func(r *http.Request) (string, string)) *RateLimiter {
	rl := &RateLimiter{
		config:       cfg,
		log:          log,
		keyExtractor: extractor,
		cleanup:      time.NewTicker(10 * time.Minute),
	}
	go rl.startCleanup()
	return rl
}

func (rl *RateLimiter) startCleanup() {
	for range rl.cleanup.C {
		rl.limiters.Range(func(key, value any) bool {
			rl.limiters.Delete(key)
			return true
		})
	}
}

func (rl *RateLimiter) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, operation := rl.keyExtractor(r)

		limitConfig := rl.config.Default
		if opConfig, ok := rl.config.Operations[operation]; ok {
			limitConfig = opConfig
		}

		// Key by IP and Operation to have per-endpoint limits per user
		key := fmt.Sprintf("%s:%s", ip, operation)

		limiter, _ := rl.limiters.LoadOrStore(key, rate.NewLimiter(rate.Limit(limitConfig.Rate), limitConfig.Burst))
		if !limiter.(*rate.Limiter).Allow() {
			rl.log.Warn(r.Context(), "rate limit exceeded",
				logger.WithField("ip", ip),
				logger.WithField("operation", operation))
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("Rate limit exceeded"))
			return
		}

		next(w, r)
	}
}

func ExtractIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		if ip == "" {
			ip = r.RemoteAddr
		}
	} else {
		if strings.Contains(ip, ",") {
			ip = strings.TrimSpace(strings.Split(ip, ",")[0])
		}
	}
	return ip
}

func DefaultRESTExtractor(r *http.Request) (string, string) {
	return ExtractIP(r), fmt.Sprintf("%s:%s", r.Method, r.URL.Path)
}
