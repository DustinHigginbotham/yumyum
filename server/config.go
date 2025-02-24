package main

import "os"

type Config struct {
	systemPrompt string
	accessToken  string
	port         string
}

func loadConfig() Config {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8666"
	}

	return Config{
		systemPrompt: os.Getenv("PROMPT"),
		accessToken:  os.Getenv("ACCESS_TOKEN"),
		port:         port,
	}
}
