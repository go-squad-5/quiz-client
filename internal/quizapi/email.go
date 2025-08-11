package quizapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetEmailReportResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func (q *QuizAPI) GetEmailReport(sessionID string) (string, error) {
	reqUrl := fmt.Sprintf(q.endpoints.getEmailReport, sessionID)

	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get email report: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		var errorResp GetEmailReportResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return "", fmt.Errorf("failed to parse error response body: %w", err)
		}
		return "", fmt.Errorf("failed to get email report, status code: %d, message: %s", errorResp.StatusCode, errorResp.Message)
	}

	return "Email report request accepted", nil
}
