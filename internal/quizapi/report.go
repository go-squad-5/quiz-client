package quizapi

import (
	"encoding/base64"
	// "encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GetReportResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		DocumentID    string `json:"documentId"`
		FileName      string `json:"fileName"`
		DownloadURL   string `json:"downloadUrl"`
		ExpiresAt     string `json:"expiresAt"`
		ContentBase64 string `json:"contentBase64"`
	} `json:"data"`
}

func (q *QuizAPI) GetReport(sessionID, userId string) (string, error) {
  // prepare request
	reqUrl := fmt.Sprintf("%s?sessionId=%s&userId=%s", q.endpoints.getReport, sessionID, userId)
	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get report: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get report, status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	// if binary stream response:
	reportPath, err := parseBinaryResponse(resp, sessionID)
	if err != nil {
		return "", fmt.Errorf("Error handling response: %w", err)
	}

	// if json - base64 response:
	// var reportResp GetReportResponse
	// if err := json.NewDecoder(resp.Body).Decode(&reportResp); err != nil {
	// 	return "", fmt.Errorf("failed to parse response body: %w", err)
	// }
	//
	//  reportPath, err := parseBase64Response(reportResp.Data.ContentBase64, sessionID)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to parse base64 response: %w", err)
	// }

	return reportPath, nil
}

func parseBase64Response(base64String, sessionId string) (string, error) {
	filePath := fmt.Sprintf("./tmp/%s/report.pdf", sessionId)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	binaryData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}
	if _, err := file.Write(binaryData); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

func parseBinaryResponse(resp *http.Response, sessionId string) (string, error) {
	filePath := fmt.Sprintf("./tmp/%s/report.pdf", sessionId)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	defer file.Close()
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}
