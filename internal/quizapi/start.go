package quizapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type StartQuizRequest struct {
	SessionID string `json:"ssid"`
}

type StartQuizResponse struct {
	SessionID string     `json:"session_id"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID       string   `json:"que_id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

func (q *QuizAPI) StartQuiz(sessionId string) ([]Question, error) {
	// Prepare Request
	req := StartQuizRequest{
		SessionID: sessionId,
	}
	if req.SessionID == "" {
		return nil, fmt.Errorf("SessionID is required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Send Request
	resp, err := q.client.Post(
		q.endpoints.startQuiz,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to start quiz, status code: %d", resp.StatusCode)
	}

	// Decode Response
	defer resp.Body.Close()
	var response StartQuizResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if response.SessionID == "" || len(response.Questions) == 0 {
		return nil, fmt.Errorf("invalid response: session_id or questions are empty")
	}

	return response.Questions, nil
}
