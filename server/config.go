package main

import (
	"encoding/base64"
	"os"
)

type Config struct {
	systemPrompt string
	accessToken  string
	port         string
	frontendURL  string
	redisURL     string
}

func loadConfig() Config {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8666"
	}

	systemPrompt := os.Getenv("PROMPT")
	if systemPrompt != "" {
		systemPrompt = decodePrompt(systemPrompt)
	}

	return Config{
		systemPrompt: systemPrompt,
		accessToken:  os.Getenv("ACCESS_TOKEN"),
		port:         port,
		frontendURL:  os.Getenv("FRONTEND_URL"),
		redisURL:     os.Getenv("REDIS_URL"),
	}
}

func decodePrompt(prompt string) string {

	decoded, err := base64.StdEncoding.DecodeString(prompt)
	if err != nil {
		return ""
	}

	return string(decoded)
}
