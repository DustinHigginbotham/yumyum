package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OpenAIMessage struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func createStream(w http.ResponseWriter) (http.Flusher, error) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming not supported")
	}

	return flusher, nil
}

// streamResponse is a helper function that sets up the event stream, converts the data from OpenAI,
// and then streams it to the client
func streamResponse(w http.ResponseWriter, resp *http.Response) {

	flusher, err := createStream(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// OpenAI streams responses in the format: "data: { JSON }"
		if !bytes.HasPrefix([]byte(line), []byte("data: ")) {
			continue // Skip any unexpected lines
		}

		jsonPart := line[6:] // Strip "data: " prefix

		// Check for end of stream
		if jsonPart == "[DONE]" {
			break
		}

		var msg OpenAIMessage
		err := json.Unmarshal([]byte(jsonPart), &msg)
		if err != nil {
			fmt.Println(fmt.Errorf("error unmarshalling OpenAI response: %w", err))
			continue
		}

		if len(msg.Choices) > 0 && msg.Choices[0].Delta.Content != "" {

			// Newlines kept getting lost somewhere so adding them back in as [NEWLINE]
			escapedContent := strings.ReplaceAll(msg.Choices[0].Delta.Content, "\n", "[NEWLINE]")

			fmt.Fprintf(w, "data: %s\n\n", escapedContent)
			flusher.Flush()
		}
	}

	// End of stream
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}
