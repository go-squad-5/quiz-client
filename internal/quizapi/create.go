package quizapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CreateSessionAPIRequest struct {
	Email string `json:"email"`
	Topic string `json:"topic"`
}

type CreateSessionAPIResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

func (q *QuizAPI) CreateSession(email, topic string) (string, error) {
	if err := validateCreateSessionInputs(email, topic); err != nil {
		return "", err
	}

	body, err := buildCreateSessionAPIRequest(email, topic)
	if err != nil {
		return "", err
	}

	// send the request
	resp, err := q.client.Post(
		q.endpoints.createSession,
		"application/json",
		body,
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err = validateCreateSessionAPIResponseStatus(resp); err != nil {
		return "", err
	}

	return parseCreateSessionAPIResponse(resp.Body)
}

func isValidEmail(email string) bool {
	if i := strings.Index(email, "@"); len(email) < 3 || i < 1 || (i != -1 && !strings.Contains(email[i:], ".")) {
		return false
	}
	return true
}

func validateCreateSessionInputs(email, topic string) error {
	if email == "" || topic == "" {
		return fmt.Errorf("email and topic are required")
	}
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}
	return nil
}

func buildCreateSessionAPIRequest(email, topic string) (*bytes.Buffer, error) {
	req := CreateSessionAPIRequest{
		Email: email,
		Topic: topic,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

func validateCreateSessionAPIResponseStatus(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create session, status code: %d", resp.StatusCode)
	}
	return nil
}

func parseCreateSessionAPIResponse(body io.ReadCloser) (string, error) {
	var response CreateSessionAPIResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}
	if response.SessionID == "" {
		return "", fmt.Errorf("session ID is empty in response: %s", response.Message)
	}
	return response.SessionID, nil
}
