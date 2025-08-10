package quizapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateSessionRequest struct {
	Email string `json:"email"`
	Topic string `json:"topic"`
}

type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

func (q *QuizAPI) CreateSession(email, topic string) (string, error) {
	// prepare the request
	req := CreateSessionRequest{
		Email: email,
		Topic: topic,
	}
	if req.Email == "" || req.Topic == "" {
		return "", fmt.Errorf("email and topic are required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// send the request
	resp, err := q.client.Post(
		q.endpoints.createSession,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create session, status code: %d", resp.StatusCode)
	}

	// decode the response
	defer resp.Body.Close()
	var response CreateSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return response.SessionID, nil
}
