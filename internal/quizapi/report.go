package quizapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"fmt"
	"io"
	"net/http"
	"os"
)

type GetReportResponseData struct {
	DocumentID    string `json:"documentId"`
	FileName      string `json:"fileName"`
	DownloadURL   string `json:"downloadUrl"`
	ExpiresAt     string `json:"expiresAt"`
	ContentBase64 string `json:"contentBase64"`
}

type GetReportResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    GetReportResponseData `json:"data"`
}

type GetReportErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (q *QuizAPI) GetReport(sessionID string) (string, error) {
	reqUrl := buildGetReportAPIURL(q.endpoints.getReport, sessionID)

	resp, err := q.client.Get(reqUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get report: %w", err)
	}
	defer resp.Body.Close()

	if err := validateGetReportResponseStatus(resp); err != nil {
		errResp, err := parseGetReportErrorResponse(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to parse error response: %w", err)
		}
		return "", errors.New(errResp)
	}

	file, filePath, err := openSessionReportFile(sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to open session report file: %w", err)
	}
	defer file.Close()

	if err := saveResponseToFile(resp, file); err != nil {
		return "", err
	}

	return filePath, nil
}

func buildGetReportAPIURL(endpoint, sessionID string) string {
	return fmt.Sprintf(endpoint, sessionID)
}

func validateGetReportResponseStatus(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get report, status code: %d", resp.StatusCode)
	}
	return nil
}

func parseGetReportErrorResponse(body io.ReadCloser) (string, error) {
	var errorResp GetReportErrorResponse
	if err := json.NewDecoder(body).Decode(&errorResp); err != nil {
		return "", fmt.Errorf("failed to parse error response body: %w", err)
	}
	if errorResp.StatusCode == 0 {
		return "", fmt.Errorf("error response should have a valid statusCode")
	}
	return fmt.Sprintf("failed to get report, status code: %d, message: %s", errorResp.StatusCode, errorResp.Message), nil
}

var reportsDirPath string = "./tmp/reports"

func openSessionReportFile(sessionId string) (*os.File, string, error) {
	// create a tmp directory if it doesn't exist
	if _, err := os.Stat(reportsDirPath); os.IsNotExist(err) {
		err := os.MkdirAll(reportsDirPath, 0755)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create tmp/reports directory: %w", err)
		}
	}

	filePath := fmt.Sprintf("%s/%s_report.pdf", reportsDirPath, sessionId)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create file: %w", err)
	}

	return file, filePath, nil
}

func saveResponseToFile(resp *http.Response, file *os.File) error {
	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// NOTE: if json - base64 response (as per API docs, we have two possible response formats)
//
// respJson, err := parseJsonGetReportResponse(resp.Body)
// if err != nil {
//   return "", fmt.Errorf("failed to parse JSON response: %w", err)
// }
//
// if err := decodeAndSaveBase64Response(respJson.Data.ContentBase64, file); err != nil {
// 	 return "", fmt.Errorf("failed to parse base64 response: %w", err)
// }

func parseJsonGetReportResponse(respBody io.ReadCloser) (GetReportResponse, error) {
	var reportResp GetReportResponse
	if err := json.NewDecoder(respBody).Decode(&reportResp); err != nil {
		return GetReportResponse{}, fmt.Errorf("failed to parse response body: %w", err)
	}
	return reportResp, nil
}

func decodeAndSaveBase64Response(base64String string, file *os.File) error {
	binaryData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return fmt.Errorf("failed to decode base64 string: %w", err)
	}
	if _, err := file.Write(binaryData); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
