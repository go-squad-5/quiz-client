package quizapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_quizapi_create_IsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"valid email", "mohit@example.com", true},
		{"invalid email without @", "mohitexample.com", false},
		{"invalid email without domain", "mohit@", false},
		{"invalid email with only @", "@example.com", false},
		{"empty email", "", false},
		{"valid email with subdomain", "mohit@abc.def.com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("isValidEmail(%s) = %v, expected %v", tt.email, result, tt.expected)
			}
		})
	}
}

func Test_quizapi_create_ValidateCreateSessionInputs(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		topic         string
		expectedError bool
	}{
		{"valid inputs", "mohit@example.com", "math", false},
		{"invalid email", "invalid-email", "math", true},
		{"empty email", "", "math", true},
		{"empty topic", "mohit@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateSessionInputs(tt.email, tt.topic)
			if (err != nil) != tt.expectedError {
				t.Errorf("validateCreateSessionInputs() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func Test_quizapi_create_BuildCreateSessionAPIRequest(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		topic         string
		expectedError bool
	}{
		{"valid inputs", "mohit@example.com", "math", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := buildCreateSessionAPIRequest(tt.email, tt.topic)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error for inputs: email=%s, topic=%s", tt.email, tt.topic)
			} else {
				assert.NoError(t, err, "Did not expect an error for inputs: email=%s, topic=%s", tt.email, tt.topic)
				bodyString := body.String()
				assert.NotEmpty(t, bodyString, "Expected non-empty request body for inputs: email=%s, topic=%s", tt.email, tt.topic)
				assert.Contains(t, bodyString, tt.email, "Expected email in request body for inputs: email=%s, topic=%s", tt.email, tt.topic)
				assert.Contains(t, bodyString, tt.topic, "Expected topic in request body for inputs: email=%s, topic=%s", tt.email, tt.topic)
			}
		})
	}
}

func Test_quizapi_create_ValidateCreateSessionAPIResponseStatus(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError bool
	}{
		{"success status", 200, false},
		{"bad request", 400, true},
		{"unauthorized", 401, true},
		{"not found", 404, true},
		{"internal server error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			err := validateCreateSessionAPIResponseStatus(resp)
			if tt.expectedError {
				assert.NotNil(t, err, "Expected an error for status code %d", tt.statusCode)
			}
		})
	}
}

func Test_quizapi_create_ParseCreateSessionAPIResponse(t *testing.T) {
	tests := []struct {
		name          string
		responseBody  string
		expectedError bool
	}{
		{"valid response", `{"session_id": "12345", "message": "success"}`, false},
		{"invalid JSON", `{"session_id": 12345`, true},
		{"missing session_id", `{"message": "success"}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionID, err := parseCreateSessionAPIResponse(
				io.NopCloser(strings.NewReader(tt.responseBody)),
			)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error for response: %s", tt.responseBody)
			} else {
				assert.NoError(t, err, "Did not expect an error for response: %s", tt.responseBody)
				assert.NotEmpty(t, sessionID, "Expected non-empty session ID for response: %s", tt.responseBody)
			}
		})
	}
}

func Test_quizapi_create_CreateSession_WhenValidInputs(t *testing.T) {
	email := "example@example.com"
	topic := "math"
	quizClient := NewTestQuizAPI("http://localhost:8080", "http://localhost:8081", func(req *http.Request) *http.Response {
		assert.Equal(t, req.Method, http.MethodPost)
		assert.Equal(t, req.URL.String(), "http://localhost:8080/session/create")
		assert.Equal(t, req.Header.Get("Content-Type"), "application/json")
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err, "Failed to read request body")
		assert.Equal(t, string(body), fmt.Sprintf(`{"email":"%s","topic":"%s"}`, email, topic))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"session_id": "12345", "message": "success"}`)),
			Header:     make(http.Header),
			Request:    req,
		}
	})
	ssid, err := quizClient.CreateSession(email, topic)
	assert.NoError(t, err, "Expected no error when creating session")
	assert.Equal(t, "12345", ssid, "Expected session ID to be '12345'")
}

func Test_quizapi_create_CreateSession_WhenInValidInputs(t *testing.T) {
	email := "example@a"
	topic := "math"

	quizClient := NewTestQuizAPI("http://localhost:8080", "http://localhost:8081", func(req *http.Request) *http.Response {
		assert.Equal(t, req.Method, http.MethodPost)
		assert.Equal(t, req.URL.String(), "http://localhost:8080/session/create")
		assert.Equal(t, req.Header.Get("Content-Type"), "application/json")
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err, "Failed to read request body")
		assert.Equal(t, string(body), fmt.Sprintf(`{"email":"%s","topic":"%s"}`, email, topic))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"session_id": "12345", "message": "success"}`)),
			Header:     make(http.Header),
			Request:    req,
		}
	})

	_, err := quizClient.CreateSession(email, topic)
	assert.Error(t, err, "Expected an error when creating session with invalid email")
}
