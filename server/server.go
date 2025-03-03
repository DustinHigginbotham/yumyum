package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const userPrompt = `You are a chef who has just created a new recipe. You named it "{{.Name}}" and it's made with only these ingredients: {{.Ingredients}}. You're excited to share it with the world, but you want to add a backstory to make it more interesting. What's the backstory of your new recipe?`

type Server struct {
	config     Config
	userPrompt string
	redis      *redis.Client
}

func NewServer(config Config) *Server {
	return &Server{
		config:     config,
		userPrompt: userPrompt,
	}
}

func (s *Server) Start() {

	// connect to redis, if defined
	if s.config.redisURL != "" {

		s.redis = redis.NewClient(&redis.Options{
			Addr: s.config.redisURL,
		})

		limiter := s.rateLimiterMiddleware()

		// set up our single route
		http.Handle("GET /generate", s.corsMiddleware(limiter(s.handleGenerateBackstory())))
	} else {
		http.Handle("GET /generate", s.corsMiddleware(s.handleGenerateBackstory()))
	}

	fmt.Printf("Server running on port %s\n", s.config.port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", s.config.port), nil)
	if err != nil {
		fmt.Printf("error starting server: %v\n", err)
	}
}

func (s *Server) handleGenerateBackstory() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Pull the query parameters from the URL
		q := r.URL.Query()
		name := q.Get("name")
		ingredients := q.Get("ingredients")

		// We are going to use a Go template to replace the placeholders in the userPrompt
		// with the actual values from the query parameters
		data := map[string]string{
			"Name":        name,
			"Ingredients": ingredients,
		}

		// Create a buffer to hold the parsed prompt
		promptBuffer := new(bytes.Buffer)
		t := template.Must(template.New("recipe").Parse(strings.ReplaceAll(s.userPrompt, "\\n", "\n")))
		err := t.Execute(promptBuffer, &data)
		if err != nil {
			fmt.Printf("error executing template: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		decodedPrompt := html.UnescapeString(promptBuffer.String())

		// Actually make the call to OpenAI with the decoded prompt
		resp, err := s.callOpenAI(decodedPrompt)
		if err != nil {
			fmt.Printf("error calling OpenAI: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// doesn't happen too often, but let's check if the response is nil
		if resp == nil {
			fmt.Println("response from OpenAI was nil")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		streamResponse(w, resp)
	}
}

func (s *Server) callOpenAI(prompt string) (*http.Response, error) {
	requestBody := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "system", "content": s.config.systemPrompt}, // System prompt
			{"role": "user", "content": prompt},                  // User input
		},
		"stream": true,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request to openai: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second, // Only applies to receiving the first header
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling openAI: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("OpenAI returned non-200 status: %d", resp.StatusCode)
	}

	return resp, nil
}

func (s *Server) Shutdown() {
	if s.redis != nil {
		err := s.redis.Close()
		if err != nil {
			fmt.Printf("error closing Redis connection: %v\n", err)
		}
	}
}
