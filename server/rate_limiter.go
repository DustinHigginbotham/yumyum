package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	requestLimit  = 2
	requestPeriod = time.Minute
)

func (s *Server) rateLimiterMiddleware() func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// get the user's IP
			ip := getClientIP(r)
			key := fmt.Sprintf("rate_limit:%s", ip)

			count, err := s.redis.Get(r.Context(), key).Int()
			if err != nil && err != redis.Nil {
				fmt.Printf("Error getting result from redis for rate limiting: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if count >= requestLimit {
				// Since we're using Server-Sent Events, we need to flush the response

				flusher, err := createStream(w)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)

				fmt.Fprint(w, "event: error\n")
				fmt.Fprint(w, "data: {\"error\": \"rate limit exceeded\"}\n\n")

				flusher.Flush()

				return
			}

			pipe := s.redis.TxPipeline()
			pipe.Incr(r.Context(), key)
			pipe.Expire(r.Context(), key, requestPeriod)
			_, err = pipe.Exec(r.Context())
			if err != nil {
				fmt.Printf("Error incrementing rate limit in redis: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)

		})
	}

}

func getClientIP(r *http.Request) string {
	var ip string

	// Check common headers set by proxies/load balancers
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can have multiple IPs, we take the first one
		ips := strings.Split(ip, ",")
		ip = strings.TrimSpace(ips[0])
	}

	if ip == "" {
		ip = r.Header.Get("CF-Connecting-IP") // Cloudflare specific header
	}

	if ip == "" {
		ip = r.Header.Get("X-Real-IP") // Some proxies use this instead
	}

	if ip == "" {
		// Fallback: Extract only the IP from RemoteAddr (stripping the port)
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("Error parsing IP from RemoteAddr:", err)
			return r.RemoteAddr // If parsing fails, return raw RemoteAddr
		}
		ip = host
	}

	fmt.Println("Resolved Client IP:", ip) // Debugging line

	return ip
}
