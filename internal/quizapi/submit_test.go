package quizapi

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_quizapi_submit_ValidateSubmitQuizInputs(t *testing.T) {
	tests := []struct {
		name          string
		sessionId     string
		answers       []Answer
		expectedError bool
	}{
		{"valid inputs", "session123", []Answer{{QuestionID: "ques1", Answer: "answer1"}}, false},
		{"missing session ID", "", []Answer{{QuestionID: "ques1", Answer: "answer1"}}, true},
		{"empty answers", "session123", []Answer{}, true},
		{"missing question ID", "session123", []Answer{{QuestionID: "", Answer: "answer1"}}, true},
		{"missing answer", "session123", []Answer{{QuestionID: "ques1", Answer: ""}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSubmitQuizInputs(tt.sessionId, tt.answers)
			if tt.expectedError {
				assert.NotNil(t, err, "Expected an error for test case: %s", tt.name)
			} else {
				assert.Nil(t, err, "Did not expect an error for test case: %s", tt.name)
			}
		})
	}
}

func Test_quizapi_submit_BuildSubmitQuizAPIRequestBody(t *testing.T) {
	sessionId := "session123"
	answers := []Answer{
		{QuestionID: "ques1", Answer: "answer1"},
		{QuestionID: "ques2", Answer: "answer2"},
	}
	expectedBody := `{"session_id":"session123","answers":[{"ques_id":"ques1","answer":"answer1"},{"ques_id":"ques2","answer":"answer2"}]}`

	body, err := buildSubmitQuizAPIRequestBody(sessionId, answers)

	require.NoError(t, err, "Expected no error while building request body")

	actualBody := body.String()
	assert.Equal(t, expectedBody, actualBody, "Request body should match expected JSON format")
}

func Test_quizapi_report_ValidateSubmitQuizAPIStatus(t *testing.T) {
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
			err := validateSubmitQuizAPIStatus(resp)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error for status code %d", tt.statusCode)
				assert.Contains(t, err.Error(), strconv.Itoa(tt.statusCode), "Expected error message to contain status code %d", tt.statusCode)
			} else {
				assert.NoError(t, err, "Did not expect an error for status code %d", tt.statusCode)
			}
		})
	}
}

func Test_quizapi_submit_ParseSubmitQuizAPIResponse(t *testing.T) {
	tests := []struct {
		name          string
		responseBody  string
		expectedScore int
		expectedError bool
	}{
		{"valid response", `{"score": 85}`, 85, false},
		{"zero score", `{"score": 0}`, 0, false},
		{"invalid response", `{"score": -1}`, 0, true},
		{"invalid response type", `{"score": "eighty-five"}`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tt.responseBody))
			score, err := parseSubmitQuizAPIResponse(body)

			if tt.expectedError {
				assert.Error(t, err, "Expected an error for test case: %s", tt.name)
				assert.Equal(t, 0, score, "Expected score to be 0 for test case: %s", tt.name)
			} else {
				assert.NoError(t, err, "Did not expect an error for test case: %s", tt.name)
				assert.Equal(t, tt.expectedScore, score, "Expected score to match for test case: %s", tt.name)
			}
		})
	}
}

func Test_quizapi_submit_SubmitQuiz_WhenValidInputs(t *testing.T) {
	sessionId := "session123"
	answers := []Answer{
		{QuestionID: "ques1", Answer: "answer1"},
		{QuestionID: "ques2", Answer: "answer2"},
	}
	expectedBody := `{"session_id":"session123","answers":[{"ques_id":"ques1","answer":"answer1"},{"ques_id":"ques2","answer":"answer2"}]}`

	expectedScore := 85
	q := NewTestQuizAPI("http://example.com", "http://report.example.com", func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method, "Expected POST method for SubmitQuiz")
		assert.Equal(t, "http://example.com/quiz/submit", req.URL.String(), "Expected correct URL for SubmitQuiz")
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "Expected Content-Type to be application/json")
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		assert.Equal(t, expectedBody, string(body), "Request body should match expected JSON format")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"score": ` + strconv.Itoa(expectedScore) + `}`)),
			Header:     make(http.Header),
		}
	})

	score, err := q.SubmitQuiz(sessionId, answers)

	assert.NoError(t, err, "Expected no error while submitting quiz")
	assert.Equal(t, expectedScore, score, "Expected score to match the expected value")
}

func Test_quizapi_submit_SubmitQuiz_WhenInvalidInputs(t *testing.T) {
	sessionId := ""
	answers := []Answer{
		{QuestionID: "ques1", Answer: "answer1"},
		{QuestionID: "ques2", Answer: "answer2"},
	}
	expectedBody := `{"session_id":"session123","answers":[{"ques_id":"ques1","answer":"answer1"},{"ques_id":"ques2","answer":"answer2"}]}`

	expectedScore := 85
	q := NewTestQuizAPI("http://example.com", "http://report.example.com", func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method, "Expected POST method for SubmitQuiz")
		assert.Equal(t, "http://example.com/quiz/submit", req.URL.String(), "Expected correct URL for SubmitQuiz")
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "Expected Content-Type to be application/json")
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		assert.Equal(t, expectedBody, string(body), "Request body should match expected JSON format")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"score": ` + strconv.Itoa(expectedScore) + `}`)),
			Header:     make(http.Header),
		}
	})

	score, err := q.SubmitQuiz(sessionId, answers)

	require.Error(t, err, "Expected no error while submitting quiz")
	assert.Equal(t, 0, score, "Expected score to match the zero value for error cases")
	assert.Contains(t, err.Error(), "SessionID is required")
}

func Test_quizapi_submit_SubmitQuiz_WhenInvalidResponse(t *testing.T) {
	sessionId := "session123"
	answers := []Answer{
		{QuestionID: "ques1", Answer: "answer1"},
		{QuestionID: "ques2", Answer: "answer2"},
	}
	expectedBody := `{"session_id":"session123","answers":[{"ques_id":"ques1","answer":"answer1"},{"ques_id":"ques2","answer":"answer2"}]}`

	expectedScore := 0
	q := NewTestQuizAPI("http://example.com", "http://report.example.com", func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method, "Expected POST method for SubmitQuiz")
		assert.Equal(t, "http://example.com/quiz/submit", req.URL.String(), "Expected correct URL for SubmitQuiz")
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "Expected Content-Type to be application/json")
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		assert.Equal(t, expectedBody, string(body), "Request body should match expected JSON format")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"score": ` + strconv.Itoa(expectedScore))),
			Header:     make(http.Header),
		}
	})

	score, err := q.SubmitQuiz(sessionId, answers)

	require.Error(t, err, "Expected no error while submitting quiz")
	assert.Equal(t, 0, score, "Expected score to match the zero value for error cases")
	assert.Contains(t, err.Error(), "failed to decode", "Expected error message to indicate decoding failure - failed to decode")
	assert.Contains(t, err.Error(), "unexpected EOF", "Expected error message to indicate correct message - unexpected EOF")
}

func Test_quizapi_submit_SubmitQuiz_WhenErrorStatusCode(t *testing.T) {
	sessionId := "session123"
	answers := []Answer{
		{QuestionID: "ques1", Answer: "answer1"},
		{QuestionID: "ques2", Answer: "answer2"},
	}
	expectedBody := `{"session_id":"session123","answers":[{"ques_id":"ques1","answer":"answer1"},{"ques_id":"ques2","answer":"answer2"}]}`

	expectedScore := 0
	q := NewTestQuizAPI("http://example.com", "http://report.example.com", func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method, "Expected POST method for SubmitQuiz")
		assert.Equal(t, "http://example.com/quiz/submit", req.URL.String(), "Expected correct URL for SubmitQuiz")
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "Expected Content-Type to be application/json")
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		assert.Equal(t, expectedBody, string(body), "Request body should match expected JSON format")
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader(`{"score": ` + strconv.Itoa(expectedScore) + `}`)),
			Header:     make(http.Header),
		}
	})

	score, err := q.SubmitQuiz(sessionId, answers)

	require.Error(t, err, "Expected error while submitting error quiz")
	assert.Equal(t, 0, score, "Expected score to match the zero value for error cases")
	assert.Contains(t, err.Error(), "401", "Expected error message to indicate the status code 401")
}

func Test_quizapi_submit_SubmitQuiz_WhenNetworkError(t *testing.T) {
	sessionId := "session123"
	answers := []Answer{
		{QuestionID: "ques1", Answer: "answer1"},
		{QuestionID: "ques2", Answer: "answer2"},
	}
	q := NewTestQuizAPI("http://example.com", "http://report.example.com", func(req *http.Request) *http.Response {
		return nil
	})
	score, err := q.SubmitQuiz(sessionId, answers)
	assert.Error(t, err, "Expected an error when network request fails")
	assert.Equal(t, 0, score, "Expected score to be 0 when there is a network error")
}
