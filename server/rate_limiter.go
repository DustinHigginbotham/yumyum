package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

const (
	requestLimit  = 1
	requestPeriod = time.Minute
)

func (s *Server) rateLimiterMiddleware() func(next http.Handler) http.Handler {

	rateLimiter := redis_rate.NewLimiter(s.redis)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// get the user's IP
			ip := getClientIP(r)

			res, err := rateLimiter.Allow(r.Context(), ip, redis_rate.Limit{Rate: requestLimit, Period: requestPeriod})

			// something happened when checking the rate limit
			if err != nil {
				fmt.Printf("error checking rate limit: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "error checking rate limit"}`))
				return
			}

			if res.Allowed == 0 {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "rate limit exceeded"}`))
				return
			}

			next.ServeHTTP(w, r)

		})
	}

}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	ips := strings.Split(ip, ",") // In case of multiple proxies
	return strings.TrimSpace(ips[0])
}
