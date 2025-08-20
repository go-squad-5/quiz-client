package app

import (
	"errors"
	"testing"
	"time"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_app_session_NewSession(t *testing.T) {
	email, topic := "test@example.com", "math"

	session := NewSession(email, topic, nil)

	require.NotNil(t, session, "NewSession should return a non-nil Session")
	assert.Equal(t, email, session.Email, "Session Email should match the input email")
	assert.Equal(t, topic, session.Topic, "Session Topic should match the input topic")
	assert.Equal(t, STATUS_STARTED, session.Status, "Session Status should be 'started' by default")
	assert.Empty(t, session.ID, "Session ID should be empty initially")
	assert.NotEmpty(t, session.StartTime, "Session StartTime should be set")
	assert.Nil(t, session.APIsTimeTaken, "Session APIsTimeTaken should be nil as passed")
}

func Test_app_session_NewSession_SetSession(t *testing.T) {
	session := NewSession("", "", nil)
	session.SetSession("sessionID")
	require.Equal(t, "sessionID", session.ID, "Session ID should be set correctly")
}

func Test_app_session_SetStatus(t *testing.T) {
	session := NewSession("", "", nil)
	session.SetStatus(STATUS_COMPLETED)
	require.Equal(t, STATUS_COMPLETED, session.Status, "Session Status should be set to 'completed'")
	session.SetStatus(STATUS_FAILED)
	require.Equal(t, STATUS_FAILED, session.Status, "Session Status should be set to 'failed'")
}

func Test_app_session_SetScore(t *testing.T) {
	session := NewSession("", "", nil)
	session.SetScore(100)
	require.Equal(t, 100, session.Score, "Session Score should be set to 100")
}

func Test_app_session_SetReport(t *testing.T) {
	session := NewSession("", "", nil)
	report := "This is a test report"
	session.SetReport(report)
	require.Equal(t, report, session.Report, "Session Report should match the input report")
}

func Test_app_session_SetAnswers(t *testing.T) {
	session := NewSession("", "", nil)
	answers := []quizapi.Answer{
		{QuestionID: "q1", Answer: "a1"},
		{QuestionID: "q2", Answer: "a2"},
	}
	session.SetAnswers(answers)
	require.Equal(t, answers, session.Answers, "Session Answers should match the input answers")
}

func Test_app_session_SetError(t *testing.T) {
	session := NewSession("", "", nil)
	err := errors.New("Test error message")
	session.SetError(err)
	require.NotNil(t, session.Error, "Session Error should not be nil")
	assert.Equal(t, err.Error(), session.Error.Error(), "Session Error message should match the input error message")
}

func Test_app_session_SetStartTime(t *testing.T) {
	session := NewSession("", "", nil)
	startTime := int64(1633072800000) // Example timestamp
	session.SetStartTime(time.UnixMilli(startTime))
	require.Equal(t, startTime, session.StartTime, "Session StartTime should be set correctly")
}

func Test_app_session_SetEndTime(t *testing.T) {
	session := NewSession("", "", nil)
	endTime := int64(1633076400000) // Example timestamp
	session.SetEndTime(time.UnixMilli(endTime))
	require.Equal(t, endTime, session.EndTime, "Session EndTime should be set correctly")
}

func Test_app_session_NewAPIsTimeTaken(t *testing.T) {
	aPIsTimeTaken := NewAPIsTimeTaken()
	require.NotNil(t, aPIsTimeTaken, "NewAPIsTimeTaken should return a non-nil APIsTimeTaken")
	assert.Equal(t, int64(0), aPIsTimeTaken.SessionCreation, "APIsTimeTaken Session Creation should be initialized to 0")
	assert.Equal(t, int64(0), aPIsTimeTaken.StartQuiz, "APIsTimeTaken Start Quiz should be initialized to 0")
	assert.Equal(t, int64(0), aPIsTimeTaken.SubmitQuiz, "APIsTimeTaken Submit Quiz should be initialized to 0")
	assert.Equal(t, int64(0), aPIsTimeTaken.ReportAPI, "APIsTimeTaken ReportAPI should be initialized to 0")
	assert.Equal(t, int64(0), aPIsTimeTaken.EmailAPI, "APIsTimeTaken Email API should be initialized to 0")
}

func Test_app_session_SetSessionCreationTime(t *testing.T) {
	aPIsTimeTaken := NewAPIsTimeTaken()
	val := int64(100)
	aPIsTimeTaken.SetSessionCreationTime(val)
	assert.Equal(t, val, aPIsTimeTaken.SessionCreation, "APIsTimeTaken Session Creation should be set to 100")
}

func Test_app_session_SetStartQuizTime(t *testing.T) {
	aPIsTimeTaken := NewAPIsTimeTaken()
	val := int64(200)
	aPIsTimeTaken.SetStartQuizTime(val)
	assert.Equal(t, val, aPIsTimeTaken.StartQuiz, "APIsTimeTaken Start Quiz should be set to 200")
}

func Test_app_session_SetSubmitQuizTime(t *testing.T) {
	aPIsTimeTaken := NewAPIsTimeTaken()
	val := int64(300)
	aPIsTimeTaken.SetSubmitQuizTime(val)
	assert.Equal(t, val, aPIsTimeTaken.SubmitQuiz, "APIsTimeTaken Submit Quiz should be set to 300")
}

func Test_app_session_SetReportAPITime(t *testing.T) {
	aPIsTimeTaken := NewAPIsTimeTaken()
	val := int64(400)
	aPIsTimeTaken.SetReportAPITime(val)
	assert.Equal(t, val, aPIsTimeTaken.ReportAPI, "APIsTimeTaken Report API should be set to 400")
}

func Test_app_session_SetEmailAPITime(t *testing.T) {
	aPIsTimeTaken := NewAPIsTimeTaken()
	val := int64(500)
	aPIsTimeTaken.SetEmailAPITime(val)
	assert.Equal(t, val, aPIsTimeTaken.EmailAPI, "APIsTimeTaken Email API should be set to 500")
}
