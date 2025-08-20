package quizapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_quizapi_NewQuizAPI(t *testing.T) {
	baseUrl := "http://localhost:3000"
	reportServerUrl := "http://localhost:3001"
	q := NewQuizAPI(baseUrl, reportServerUrl)

	if !assert.NotNil(t, q, "Expected return value to be not nil, but got nil") {
		return
	}

	assert.Equal(t, 60*time.Second, q.client.Timeout, "Expected client timeout to be 60 seconds, but got %v", q.client.Timeout)
	assert.Equal(t, baseUrl+"/session/create", q.endpoints.createSession, "Expected createSession endpoint to be %s, but got %s", baseUrl+"/session/create", q.endpoints.createSession)
	assert.Equal(t, baseUrl+"/quiz/start", q.endpoints.startQuiz, "Expected startQuiz endpoint to be %s, but got %s", baseUrl+"/quiz/start", q.endpoints.startQuiz)
	assert.Equal(t, baseUrl+"/quiz/submit", q.endpoints.submitQuiz, "Expected submitQuiz endpoint to be %s, but got %s", baseUrl+"/quiz/submit", q.endpoints.submitQuiz)
	assert.Equal(t, reportServerUrl+"/sessions/%s/report", q.endpoints.getReport, "Expected getReport endpoint to be %s, but got %s", reportServerUrl+"/sessions/%s/report", q.endpoints.getReport)
	assert.Equal(t, reportServerUrl+"/sessions/%s/email-report", q.endpoints.getEmailReport, "Expected getEmailReport endpoint to be %s, but got %s", reportServerUrl+"/sessions/%s/email-report", q.endpoints.getEmailReport)
}
