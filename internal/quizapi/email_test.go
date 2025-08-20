package quizapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_quizapi_email_BuildGetEmailReportAPIURL(t *testing.T) {
	sessionID := "12345"
	expectedURL := fmt.Sprintf("http://example.com/sessions/%s/email-report", sessionID)
	actualURL := buildGetEmailReportAPIURL("http://example.com/sessions/%s/email-report", sessionID)
	assert.Equal(t, expectedURL, actualURL, "URLs should match")
}

func Test_quizapi_email_ParseGetEmailReportErrorResponse(t *testing.T) {
	errorResponse := `{"message": "Email report not found", "statusCode": 404}`
	body := io.NopCloser(strings.NewReader(errorResponse))

	errMsg, err := parseGetEmailReportErrorResponse(body)

	assert.Nil(t, err, "Expected no error to be returned")
	assert.Contains(t, errMsg, "Email report not found", "Expected error message to contain 'Email report not found'")
	assert.Contains(t, errMsg, "status code: 404", "Expected error message to contain 'status code: 404'")
}

func Test_quizapi_email_ValidateGetEmailReportResponseStatus(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError bool
	}{
		{"accepted status", 202, false},
		{"bad request", 400, true},
		{"unauthorized", 401, true},
		{"not found", 404, true},
		{"internal server error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			err := validateGetEmailReportResponseStatus(resp)
			if tt.expectedError {
				assert.NotNil(t, err, "Expected an error for status code %d", tt.statusCode)
			} else {
				assert.Nil(t, err, "Did not expect an error for status code %d", tt.statusCode)
			}
		})
	}
}

func Test_quizapi_email_GetEmailReport_WhenValidSessionID(t *testing.T) {
	sessionID := "12345"
	expectedResponse := "Email report request accepted"

	q := NewTestQuizAPI("http://example.com", "http://reportserver.com", func(req *http.Request) *http.Response {
		assert.Equal(t, http.MethodPost, req.Method, "Expected POST method for GetEmailReport")
		assert.Equal(t, fmt.Sprintf("http://reportserver.com/sessions/%s/email-report", sessionID), req.URL.String(), "Expected correct URL for GetEmailReport")

		// Simulate a successful response
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(strings.NewReader(expectedResponse)),
			Header:     make(http.Header),
		}
	})

	response, err := q.GetEmailReport(sessionID)

	assert.Nil(t, err, "Expected no error when getting email report")
	assert.Equal(t, expectedResponse, response, "Expected response to match")
}

func Test_quizapi_email_GetEmailReport_WhenErrorResponse(t *testing.T) {
	sessionID := "12345"
	errorResponse := `{"message": "Email report not found", "statusCode": 404}`

	q := NewTestQuizAPI("http://example.com", "http://reportserver.com", func(req *http.Request) *http.Response {
		assert.Equal(t, http.MethodPost, req.Method, "Expected POST method for GetEmailReport")
		assert.Equal(t, fmt.Sprintf("http://reportserver.com/sessions/%s/email-report", sessionID), req.URL.String(), "Expected correct URL for GetEmailReport")

		// Simulate an error response
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(errorResponse)),
			Header:     make(http.Header),
		}
	})

	response, err := q.GetEmailReport(sessionID)

	assert.NotNil(t, err, "Expected an error when getting email report")
	assert.Empty(t, response, "Expected empty response on error")
}
