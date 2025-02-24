package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strings"
)

const userPrompt = `You are a chef who has just created a new recipe. You named it "{{.Name}}" and it's made with only these ingredients: {{.Ingredients}}. You're excited to share it with the world, but you want to add a backstory to make it more interesting. What's the backstory of your new recipe?`

type Server struct {
	config     Config
	userPrompt string
}

func NewServer(config Config) *Server {
	return &Server{
		config:     config,
		userPrompt: userPrompt,
	}
}

func (s *Server) Start() {
	// set up our single route
	http.HandleFunc("/generate", s.handleGenerateBackstory())

	fmt.Printf("Server running on port %s\n", s.config.port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", s.config.port), nil)
	if err != nil {
		fmt.Printf("error starting server: %v\n", err)
	}
}

func (s *Server) handleGenerateBackstory() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()
		name := q.Get("name")
		ingredients := q.Get("ingredients")

		data := map[string]string{
			"Name":        name,
			"Ingredients": ingredients,
		}

		promptBuffer := new(bytes.Buffer)
		t := template.Must(template.New("recipe").Parse(strings.ReplaceAll(s.userPrompt, "\\n", "\n")))
		err := t.Execute(promptBuffer, &data)
		if err != nil {
			fmt.Printf("error executing template: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		decodedPrompt := html.UnescapeString(promptBuffer.String())

		resp, err := s.callOpenAI(decodedPrompt)
		if err != nil {
			fmt.Printf("error calling OpenAI: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		allowOrigin := "*"
		if s.config.frontendURL != "" {
			allowOrigin = s.config.frontendURL
		}

		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)

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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling openAI: %w", err)
	}

	return resp, nil
}
