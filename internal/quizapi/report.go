package quizapi

import (
	"encoding/base64"
	"encoding/json"

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

type GetReportErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (q *QuizAPI) GetReport(sessionID string) (string, error) {
	// prepare request
	reqUrl := fmt.Sprintf(q.endpoints.getReport, sessionID)
	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get report: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
    if err != nil {
      return "", fmt.Errorf("failed to read error response body: %w", err)
    }
		fmt.Println("Response body:", string(body))
		var errorResp GetReportErrorResponse
    if err := json.Unmarshal(body, &errorResp); err != nil {
      return "", fmt.Errorf("failed to parse error response: %w", err)
    }
		return "", fmt.Errorf("failed to get report, status code: %d, message: %s", errorResp.StatusCode, errorResp.Message)
	}
	defer resp.Body.Close()

	// if binary stream response:
	reportPath, err := parseAndSaveBinaryResponse(resp, sessionID)
	if err != nil {
		return "", fmt.Errorf("Error handling response: %w", err)
	}

	// NOTE: if json - base64 response:
	//
	// var reportResp GetReportResponse
	// if err := json.NewDecoder(resp.Body).Decode(&reportResp); err != nil {
	// 	return "", fmt.Errorf("failed to parse response body: %w", err)
	// }
	//
	//  reportPath, err := parseAndSaveBase64Response(reportResp.Data.ContentBase64, sessionID)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to parse base64 response: %w", err)
	// }

	return reportPath, nil
}

func parseAndSaveBase64Response(base64String, sessionId string) (string, error) {
	// make directory if not exists
	// create a tmp directory if it doesn't exist
	if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
		err := os.Mkdir("./tmp", 0755)
		if err != nil {
			panic("Failed to create tmp directory: " + err.Error())
		}
	}
	if _, err := os.Stat("./tmp/reports"); os.IsNotExist(err) {
		err := os.Mkdir("./tmp/reports", 0755)
		if err != nil {
			panic("Failed to create tmp/reports directory: " + err.Error())
		}
	}
	filePath := fmt.Sprintf("./tmp/reports/%s_report.pdf", sessionId)
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

func parseAndSaveBinaryResponse(resp *http.Response, sessionId string) (string, error) {
	filePath := fmt.Sprintf("./tmp/reports/%s_report.pdf", sessionId)
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
