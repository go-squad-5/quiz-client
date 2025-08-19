package quizapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type StartQuizAPIRequest struct {
	SessionID string `json:"ssid"`
	Topic     string `json:"topic"`
}

type StartQuizAPIResponse struct {
	SessionID string     `json:"session_id"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID       string   `json:"ques_id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

func (q *QuizAPI) StartQuiz(sessionId, topic string) ([]Question, error) {
	if err := validateStartQuizInputs(sessionId, topic); err != nil {
		return nil, err
	}

	body, err := buildStartQuizRequestBody(sessionId, topic)
	if err != nil {
		return nil, err
	}

	// Send Request
	resp, err := q.client.Post(
		q.endpoints.startQuiz,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateStartQuizAPIStatus(resp); err != nil {
		return nil, err
	}

	return parseStartQuizResponse(resp.Body)
}

func validateStartQuizInputs(sessionId, topic string) error {
	if sessionId == "" {
		return fmt.Errorf("SessionID is required")
	}
	if topic == "" {
		return fmt.Errorf("Topic is required")
	}

	return nil
}

func buildStartQuizRequestBody(sessionId, topic string) (*bytes.Buffer, error) {
	req := StartQuizAPIRequest{
		SessionID: sessionId,
		Topic:     topic,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

func validateStartQuizAPIStatus(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start quiz, status code: %d", resp.StatusCode)
	}
	return nil
}

func parseStartQuizResponse(body io.ReadCloser) ([]Question, error) {
	var response StartQuizAPIResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	if response.SessionID == "" || len(response.Questions) == 0 {
		return nil, fmt.Errorf("invalid response: session_id or questions are empty")
	}
	return response.Questions, nil
}
