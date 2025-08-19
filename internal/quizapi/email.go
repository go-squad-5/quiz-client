package quizapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetEmailReportErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func (q *QuizAPI) GetEmailReport(sessionID string) (string, error) {
	reqUrl := buildGetEmailReportAPIURL(q.endpoints.getEmailReport, sessionID)

	// send the request
	resp, err := q.client.Post(reqUrl, "application/json", nil)
	if err != nil {
		return "", fmt.Errorf("failed to send request to get email report: %w", err)
	}
	defer resp.Body.Close()

	if err := validateGetEmailReportResponseStatus(resp); err != nil {
		err = parseGetEmailReportErrorResponse(resp.Body)
		return "", err
	}

	return "Email report request accepted", nil
}

func buildGetEmailReportAPIURL(endpoint, sessionID string) string {
	return fmt.Sprintf(endpoint, sessionID)
}

func parseGetEmailReportErrorResponse(body io.ReadCloser) error {
	var errorResp GetEmailReportErrorResponse
	if err := json.NewDecoder(body).Decode(&errorResp); err != nil {
		return fmt.Errorf("failed to parse error response body: %w", err)
	}
	return fmt.Errorf("failed to get email report, status code: %d, message: %s", errorResp.StatusCode, errorResp.Message)
}

func validateGetEmailReportResponseStatus(resp *http.Response) error {
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("not accepted")
	}
	return nil
}
