package quizapi

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_quizapi_start_ValidStartQuizInputs(t *testing.T) {
	tests := []struct {
		name          string
		sessionId     string
		topic         string
		expectedError bool
	}{
		{"valid inputs", "1234", "math", false},
		{"empty sessionId", "", "science", true},
		{"empty topic", "1234", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStartQuizInputs(tt.sessionId, tt.topic)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error for inputs: sessionId=%s, topic=%s", tt.sessionId, tt.topic)
			} else {
				assert.NoError(t, err, "Did not expect an error for inputs: sessionId=%s, topic=%s", tt.sessionId, tt.topic)
			}
		})
	}
}

func Test_quizapi_start_BuildStartQuizRequestBody(t *testing.T) {
	tests := []struct {
		name          string
		sessionId     string
		topic         string
		expectedBody  string
		expectedError bool
	}{
		{"valid inputs", "1234", "math", `{"ssid":"1234","topic":"math"}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := buildStartQuizRequestBody(tt.sessionId, tt.topic)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error for inputs: sessionId=%s, topic=%s", tt.sessionId, tt.topic)
			} else {
				assert.NoError(t, err, "Did not expect an error for inputs: sessionId=%s, topic=%s", tt.sessionId, tt.topic)
				assert.NotNil(t, body, "Expected body buffer to be non-nil for inputs: sessionId=%s, topic=%s", tt.sessionId, tt.topic)
				bodyStr := body.String()
				assert.JSONEq(t, tt.expectedBody, bodyStr, "Expected body to match for inputs: sessionId=%s, topic=%s", tt.sessionId, tt.topic)
			}
		})
	}
}

func Test_quizapi_report_ValidateStartQuizAPIStatus(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError bool
	}{
		{"success status", 200, false},
		{"bad request", 400, true},
		{"unauthorized", 401, true},
		{"not found", 404, true},
		{"internal server error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			err := validateStartQuizAPIStatus(resp)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error for status code %d", tt.statusCode)
				assert.Contains(t, err.Error(), strconv.Itoa(tt.statusCode), "Expected error message to contain status code %d", tt.statusCode)
			} else {
				assert.NoError(t, err, "Did not expect an error for status code %d", tt.statusCode)
			}
		})
	}
}

func Test_quizapi_start_ParseStartQuizResponse(t *testing.T) {
	tests := []struct {
		name          string
		responseBody  string
		sessionId     string
		questions     []Question
		expectedError bool
	}{
		{
			"valid response",
			`{"session_id":"1234","questions":[{"ques_id":"q1","question":"What is 2+2?","options":["3","4","5"]}]}`,
			"1234",
			[]Question{{ID: "q1", Question: "What is 2+2?", Options: []string{"3", "4", "5"}}},
			false,
		},
		{
			"invalid JSON",
			`{"session_id":"1234","questions":[{"ques_id":"q1","question":"What is 2+2?","options":["3","4","5"]}`,
			"1234",
			nil,
			true,
		},
		{
			"empty questions",
			`{"session_id":"1234","questions":[]}`,
			"1234",
			[]Question{},
			true,
		},
		{
			"missing session_id",
			`{"questions":[{"ques_id":"q1","question":"What is 2+2?","options":["3","4","5"]}]}`,
			"",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tt.responseBody))
			questions, err := parseStartQuizResponse(body)
			if tt.expectedError {
				require.Error(t, err, "Expected an error for response: %s", tt.responseBody)
				return
			}
			require.NoError(t, err, "Did not expect an error for response: %s", tt.responseBody)
			assert.Equal(t, len(tt.questions), len(questions), "Expected number of questions to match")
			assert.ElementsMatch(t, tt.questions, questions, "Expected questions list to match")
		})
	}
}

func Test_quizapi_start_StartQuiz_WhenValidInputs(t *testing.T) {
	sessionId := "1234"
	topic := "math"
	expectedResponse := `{"session_id":"1234","questions":[{"ques_id":"q1","question":"What is 2+2?","options":["3","4","5"]}]}`
	questions := []Question{
		{ID: "q1", Question: "What is 2+2?", Options: []string{"3", "4", "5"}},
	}
	baseUrl := "http://example.com"
	reportServerUrl := "http://reportserver.com"

	q := NewTestQuizAPI(baseUrl, reportServerUrl, func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method, "Expected POST method for StartQuiz")
		expectedUrl := baseUrl + "/quiz/start"
		actualUrl := req.URL.String()
		assert.Equal(t, expectedUrl, actualUrl, "Expected URL to match for StartQuiz to be %s, got %s", expectedUrl, actualUrl)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "Expected Content-Type to be application/json")
		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err, "Expected no error reading request body")
		assert.JSONEq(t, fmt.Sprintf(`{"ssid":"%s","topic":"%s"}`, sessionId, topic), string(body), "Expected request body to match for StartQuiz")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(expectedResponse)),
			Header:     make(http.Header),
		}
	})

	response, err := q.StartQuiz(sessionId, topic)
	assert.NoError(t, err, "Expected no error when starting quiz with valid inputs")
	assert.NotNil(t, response, "Expected response to be non-nil when starting quiz with valid inputs")
	assert.Equal(t, questions, response, "Expected response questions to match the expected questions")
}

func Test_quizapi_start_StartQuiz_WhenInvalidInputs(t *testing.T) {
	sessionId := ""
	topic := "math"
	expectedResponse := `{"message":"SessionID is required","statusCode":400}`
	baseUrl := "http://example.com"
	reportServerUrl := "http://reportserver.com"

	q := NewTestQuizAPI(baseUrl, reportServerUrl, func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method, "Expected POST method for StartQuiz")
		expectedUrl := baseUrl + "/quiz/start"
		actualUrl := req.URL.String()
		assert.Equal(t, expectedUrl, actualUrl, "Expected URL to match for StartQuiz to be %s, got %s", expectedUrl, actualUrl)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "Expected Content-Type to be application/json")
		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err, "Expected no error reading request body")
		assert.JSONEq(t, fmt.Sprintf(`{"ssid":"%s","topic":"%s"}`, sessionId, topic), string(body), "Expected request body to match for StartQuiz")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(expectedResponse)),
			Header:     make(http.Header),
		}
	})

	response, err := q.StartQuiz(sessionId, topic)
	assert.Error(t, err, "Expected an error when starting quiz with invalid sessionId")
	assert.Nil(t, response, "Expected response to be nil when starting quiz with invalid inputs")
}

func Test_quizapi_start_StartQuiz_WhenErrorResponse(t *testing.T) {
	sessionId := "1234"
	topic := "math"
	expectedResponse := `{"message": "Unauthorized", "statusCode": 401}`
	baseUrl := "http://example.com"
	reportServerUrl := "http://reportserver.com"

	q := NewTestQuizAPI(baseUrl, reportServerUrl, func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method, "Expected POST method for StartQuiz")
		expectedUrl := baseUrl + "/quiz/start"
		actualUrl := req.URL.String()
		assert.Equal(t, expectedUrl, actualUrl, "Expected URL to match for StartQuiz to be %s, got %s", expectedUrl, actualUrl)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "Expected Content-Type to be application/json")
		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err, "Expected no error reading request body")
		assert.JSONEq(t, fmt.Sprintf(`{"ssid":"%s","topic":"%s"}`, sessionId, topic), string(body), "Expected request body to match for StartQuiz")
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader(expectedResponse)),
			Header:     make(http.Header),
		}
	})

	response, err := q.StartQuiz(sessionId, topic)
	assert.Error(t, err, "Expected an error when starting quiz with error response")
	assert.Nil(t, response, "Expected response to be nil when starting quiz with error response")
	assert.Contains(t, err.Error(), "401", "Expected error message to contain status code 401")
}

func Test_quizapi_start_StartQuiz_WhenInvalidSuccessResponse(t *testing.T) {
	sessionId := "1234"
	topic := "math"
	expectedResponse := `{"session_id":"1234","questions":[{"ques_id":"q1","question":"What is 2+2?","options":["3","4","5"]`
	baseUrl := "http://example.com"
	reportServerUrl := "http://reportserver.com"

	q := NewTestQuizAPI(baseUrl, reportServerUrl, func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method, "Expected POST method for StartQuiz")
		expectedUrl := baseUrl + "/quiz/start"
		actualUrl := req.URL.String()
		assert.Equal(t, expectedUrl, actualUrl, "Expected URL to match for StartQuiz to be %s, got %s", expectedUrl, actualUrl)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "Expected Content-Type to be application/json")
		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err, "Expected no error reading request body")
		assert.JSONEq(t, fmt.Sprintf(`{"ssid":"%s","topic":"%s"}`, sessionId, topic), string(body), "Expected request body to match for StartQuiz")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(expectedResponse)),
			Header:     make(http.Header),
		}
	})

	response, err := q.StartQuiz(sessionId, topic)
	assert.Error(t, err, "Expected an error when starting quiz with valid inputs with invalid response")
	assert.Nil(t, response, "Expected response to be nil when starting quiz with valid inputs with invalid response")
	assert.Contains(t, err.Error(), "failed to decode", "Expected error message to indicate decoding failure - failed to decode")
	assert.Contains(t, err.Error(), "unexpected EOF", "Expected error message to indicate correct message - unexpected EOF")
}

func Test_quizapi_start_StartQuiz_WhenNetworkError(t *testing.T) {
	sessionId := "1234"
	topic := "math"
	baseUrl := "http://example.com"
	reportServerUrl := "http://reportserver.com"

	q := NewTestQuizAPI(baseUrl, reportServerUrl, func(req *http.Request) *http.Response {
		return nil
	})

	response, err := q.StartQuiz(sessionId, topic)
	assert.Error(t, err, "Expected an error when starting quiz with network error")
	assert.Nil(t, response, "Expected response to be nil when starting quiz with network error")
}
