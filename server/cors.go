package main

import (
	"net/http"
)

func (s *Server) corsMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Open this up for cross-origin requests
		allowOrigin := "*"
		if s.config.frontendURL != "" {
			allowOrigin = s.config.frontendURL
		}

		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)

		next.ServeHTTP(w, r)
	})
}
