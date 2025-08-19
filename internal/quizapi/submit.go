package quizapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Answer struct {
	QuestionID string `json:"ques_id"`
	Answer     string `json:"answer"`
}

type SubmitQuizAPIRequest struct {
	SessionID string   `json:"session_id"`
	Answers   []Answer `json:"answers"`
}

type SubmitQuizAPIResponse struct {
	Score int `json:"score"`
}

func (q *QuizAPI) SubmitQuiz(sessionId string, answers []Answer) (int, error) {
	if err := validateSubmitQuizInputs(sessionId, answers); err != nil {
		return 0, err
	}

	body, err := buildSubmitQuizAPIRequestBody(sessionId, answers)
	if err != nil {
		return 0, fmt.Errorf("failed to build request body: %w", err)
	}

	resp, err := q.client.Post(
		q.endpoints.submitQuiz,
		"application/json",
		body,
	)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if err := validateSubmitQuizAPIStatus(resp); err != nil {
		return 0, err
	}

	return parseSubmitQuizAPIResponse(resp.Body)
}

func validateSubmitQuizInputs(sessionId string, answers []Answer) error {
	if sessionId == "" {
		return fmt.Errorf("SessionID is required")
	}
	if len(answers) == 0 {
		return fmt.Errorf("at least one answer is required")
	}
	for _, answer := range answers {
		if answer.QuestionID == "" || answer.Answer == "" {
			return fmt.Errorf("each answer must have a question ID and an answer")
		}
	}
	return nil
}

func buildSubmitQuizAPIRequestBody(sessionId string, answers []Answer) (*bytes.Buffer, error) {
	req := SubmitQuizAPIRequest{
		SessionID: sessionId,
		Answers:   answers,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	return bytes.NewBuffer(body), nil
}

func validateSubmitQuizAPIStatus(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to submit quiz, status code: %d", resp.StatusCode)
	}
	return nil
}

func parseSubmitQuizAPIResponse(body io.ReadCloser) (int, error) {
	var response SubmitQuizAPIResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}
	if response.Score < 0 {
		return 0, fmt.Errorf("invalid score received: %d", response.Score)
	}
	return response.Score, nil
}
