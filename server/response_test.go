package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DustinHigginbotham/yumyum/testdata"
)

func TestStreamResponse(t *testing.T) {

	tt := []struct {
		name          string
		data          string
		expectedLines []string
	}{
		{
			name: "Successfully parses OpenAI responses and maps to our own",
			data: testdata.MockResponseData,
			expectedLines: []string{
				"data: Grandma's Chef Boyardee is a time-honored tradition",
				"data:  passed down through generations.[NEWLINE]",
				"data: Each bite contains memories of warmth.",
				"data: [DONE]",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			mockResponse := httptest.NewRecorder()
			mockBody := io.NopCloser(strings.NewReader(tc.data))
			mockHTTPResponse := &http.Response{
				Body:       mockBody,
				StatusCode: http.StatusOK,
			}

			streamResponse(mockResponse, mockHTTPResponse)

			resultLines := readStream(strings.NewReader(mockResponse.Body.String()))

			// Ensure expected and actual have the same number of lines
			if len(resultLines) != len(tc.expectedLines) {
				t.Errorf("Expected %d lines but got %d lines", len(tc.expectedLines), len(resultLines))
			}

			for i, expected := range tc.expectedLines {
				if strings.TrimSpace(resultLines[i]) != expected {
					t.Errorf("Mismatch on line %d. Expected: %q, Got: %q", i, expected, resultLines[i])
				}
			}
		})
	}
}

func TestOutputDoesntGetWeirdOrDisturbing(t *testing.T) {
	mockResponse := httptest.NewRecorder()
	mockBody := io.NopCloser(strings.NewReader(testdata.MockResponseData2))
	mockHTTPResponse := &http.Response{
		Body:       mockBody,
		StatusCode: http.StatusOK,
	}

	streamResponse(mockResponse, mockHTTPResponse)

	resultLines := readStream(strings.NewReader(mockResponse.Body.String()))

	for _, line := range resultLines {
		if containsUnsettlingContent(line) {
			t.Errorf("Oh no, it's happening again: %q", line)
		}
	}
}
