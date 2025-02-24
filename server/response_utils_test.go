package main

import (
	"bufio"
	"io"
	"strings"
)

func readStream(response io.Reader) []string {

	// Read the response body line by line as OpenAI would stream it
	scanner := bufio.NewScanner(response)
	var resultLines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text()) // Read each line, stripping extra whitespace
		if line != "" {                           // Ignore blank lines
			resultLines = append(resultLines, line)
		}
	}

	return resultLines
}

var cursedPhrases = []string{
	"clinging to the walls",
	"the scent so familiar it aches",
	"the kitchen is empty",
	"I don't remember cooking this",
	"it has always been waiting",
}

func containsUnsettlingContent(text string) bool {
	for _, phrase := range cursedPhrases {
		if strings.Contains(strings.ToLower(text), phrase) {
			return true
		}
	}
	return false
}
