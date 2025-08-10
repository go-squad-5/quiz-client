package quizapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Answer struct {
	QuestionID string `json:"ques_id"`
	Answer     string `json:"answer"`
}

type SubmitQuizRequest struct {
	SessionID string   `json:"session_id"`
	Answers   []Answer `json:"answers"`
}

type SubmitQuizResponse struct {
	Score int `json:"score"`
}

func (q *QuizAPI) SubmitQuiz(sessionId string, answers []Answer) (int, error) {
  if sessionId == "" {
    return 0, fmt.Errorf("SessionID is required")
  }
  if len(answers) == 0 {
    return 0, fmt.Errorf("at least one answer is required")
  }

  req := SubmitQuizRequest{
    SessionID: sessionId,
    Answers:   answers,
  }

  body, err := json.Marshal(req)
  if err != nil {
    return 0, err
  }

  resp, err := q.client.Post(
    q.endpoints.submitQuiz,
    "application/json",
    bytes.NewBuffer(body),
  )
  if err != nil {
    return 0, err
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return 0, fmt.Errorf("failed to submit quiz, status code: %d", resp.StatusCode)
  }

  var response SubmitQuizResponse
  if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
    return 0, fmt.Errorf("failed to decode response: %v", err)
  }

  return response.Score, nil
}
