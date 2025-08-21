package app

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_app_results_OpenResultsFile_WhenDirExist(t *testing.T) {
	tmpDirPath = "../../tmp" // change path for the test environment
	expectedFilePath := fmt.Sprintf("%s/logs.txt", tmpDirPath)

	defer func() {
		if err := os.Remove(expectedFilePath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Error cleaning up the test file: %s", expectedFilePath)
		}
	}()

	defer func() {
		r := recover()
		require.Nilf(t, r, "Expected openResultsFile() to successfully open results file, but resulted in panic. %v", r)
	}()

	file := openResultsFile()
	defer file.Close()

	require.NotNil(t, file)
}

func Test_app_results_OpenResultsFile_WhenDirNotExist(t *testing.T) {
	tmpDirPath = "../../test" // change path for the test environment
	expectedFilePath := fmt.Sprintf("%s/logs.txt", tmpDirPath)

	defer func() {
		if err := os.RemoveAll(tmpDirPath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Error cleaning up the test files: %s", expectedFilePath)
		}
	}()

	defer func() {
		r := recover()
		require.Nilf(t, r, "Expected openResultsFile() to successfully open results file, but resulted in panic. %v", r)
	}()

	if err := os.RemoveAll(tmpDirPath); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Error while remove existing tmp dir for test")
	}

	file := openResultsFile()
	defer file.Close()

	require.NotNil(t, file)
	require.NotNil(t, file, "Expected openResultsFile() to return a file, but got nil.")
}

func Test_app_results_OpenResultsFile_WhenInvalidDirPath(t *testing.T) {
	tmpDirPath = "../..\\0/tmp"
	expectedFilePath := fmt.Sprintf("%s/logs.txt", tmpDirPath)

	defer func() {
		if err := os.Remove(expectedFilePath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Error cleaning up the test file: %s", expectedFilePath)
		}
	}()

	defer func() {
		r := recover()
		require.NotNilf(t, r, "Expected openResultsFile() to fail to open results file, but resulted in no panic. %v", r)
	}()

	if err := os.RemoveAll(tmpDirPath); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Error while removing existing tmp dir for test")
	}

	file := openResultsFile()
	defer file.Close()

	assert.Nil(t, file, "Expected openResultsFile() to return a nil file pointer, but got non-nil.")
}

func Test_app_results_GetResultLog_ValidSession(t *testing.T) {
	startTime := time.Now().UnixMilli()
	endTime := startTime + 1000
	result := &Session{
		ID:        "12345",
		UserID:    "testuser",
		Email:     "test@example.com",
		Topic:     "test-topic",
		StartTime: startTime,
		EndTime:   endTime,
		Question: []quizapi.Question{
			{
				ID:       "q1",
				Question: "What is the capital of France?",
				Options:  []string{"Paris", "London", "Berlin", "Madrid"},
			},
		},
		Answers: []quizapi.Answer{
			{
				QuestionID: "q1",
				Answer:     "Paris",
			},
		},
		Status: STATUS_COMPLETED,
		Score:  100,
		Report: "./tmp/reports/a.pdf",
		APIsTimeTaken: &APIsTimeTaken{
			SessionCreation: 500,
			StartQuiz:       300,
			SubmitQuiz:      200,
			ReportAPI:       400,
			EmailAPI:        300,
		},
	}

	logString := getResultLog(result)
	require.NotEmpty(t, logString, "Expected getResultLog() to return a non-empty string")
	assert.Contains(t, logString, result.ID, "Expected log to contain session ID")
	assert.Contains(t, logString, result.Email, "Expected log to contain email")
	assert.Contains(t, logString, result.UserID, "Expected log to contain user ID")
	assert.Contains(t, logString, fmt.Sprintf("%d", result.Score), "Expected log to contain score")
	assert.Contains(t, logString, string(result.Status), "Expected log to contain status")
	assert.Contains(t, logString, time.UnixMilli(result.StartTime).Format(time.RFC3339), "Expected log to contain start time")
	assert.Contains(t, logString, time.UnixMilli(result.EndTime).Format(time.RFC3339), "Expected log to contain end time")
	assert.Contains(t, logString, fmt.Sprintf("%d ms", result.EndTime-result.StartTime), "Expected log to contain time taken")
	assert.Contains(t, logString, fmt.Sprintf("%v", result.Answers), "Expected log to contain answers")
	assert.Contains(t, logString, result.Report, "Expected log to contain report path")
	assert.Contains(t, logString, fmt.Sprintf("Session Creation: %d ms", result.APIsTimeTaken.SessionCreation), "Expected log to contain session creation time")
	assert.Contains(t, logString, fmt.Sprintf("Start Quiz: %d ms", result.APIsTimeTaken.StartQuiz), "Expected log to contain start quiz time")
	assert.Contains(t, logString, fmt.Sprintf("Submit Quiz: %d ms", result.APIsTimeTaken.SubmitQuiz), "Expected log to contain submit quiz time")
	assert.Contains(t, logString, fmt.Sprintf("Report API: %d ms", result.APIsTimeTaken.ReportAPI), "Expected log to contain report API time")
	assert.Contains(t, logString, fmt.Sprintf("Email API: %d ms", result.APIsTimeTaken.EmailAPI), "Expected log to contain email API time")
}

func Test_app_results_GetResultLog_ErrorStartSession(t *testing.T) {
	startTime := time.Now().UnixMilli()
	endTime := startTime + 1000
	result := &Session{
		ID:        "",
		UserID:    "testuser",
		Email:     "test@example.com",
		Topic:     "test-topic",
		StartTime: startTime,
		EndTime:   endTime,
		Question: []quizapi.Question{
			{
				ID:       "q1",
				Question: "What is the capital of France?",
				Options:  []string{"Paris", "London", "Berlin", "Madrid"},
			},
		},
		Answers: []quizapi.Answer{
			{
				QuestionID: "q1",
				Answer:     "Paris",
			},
		},
		Status: STATUS_COMPLETED,
		Error:  errors.New("Test error"),
		Score:  100,
		Report: "./tmp/reports/a.pdf",
		APIsTimeTaken: &APIsTimeTaken{
			SessionCreation: 500,
			StartQuiz:       300,
			SubmitQuiz:      200,
			ReportAPI:       400,
			EmailAPI:        300,
		},
	}

	logString := getResultLog(result)
	require.NotEmpty(t, logString, "Expected getResultLog() to return a non-empty string")
	assert.Contains(t, logString, "Error while starting", "Expected log to contain Error indicator")
	assert.Contains(t, logString, result.Email, "Expected log to contain email")
	assert.Contains(t, logString, result.UserID, "Expected log to contain user ID")
	assert.Contains(t, logString, fmt.Sprintf("%d", result.Score), "Expected log to contain score")
	assert.Contains(t, logString, string(result.Status), "Expected log to contain status")
	assert.Contains(t, logString, time.UnixMilli(result.StartTime).Format(time.RFC3339), "Expected log to contain start time")
	assert.Contains(t, logString, time.UnixMilli(result.EndTime).Format(time.RFC3339), "Expected log to contain end time")
	assert.Contains(t, logString, fmt.Sprintf("%d ms", result.EndTime-result.StartTime), "Expected log to contain time taken")
	assert.Contains(t, logString, fmt.Sprintf("%v", result.Answers), "Expected log to contain answers")
	assert.Contains(t, logString, result.Report, "Expected log to contain report path")
	assert.Contains(t, logString, "Error: Test error", "Expected log to contain error message")
	assert.Contains(t, logString, fmt.Sprintf("Session Creation: %d ms", result.APIsTimeTaken.SessionCreation), "Expected log to contain session creation time")
	assert.Contains(t, logString, fmt.Sprintf("Start Quiz: %d ms", result.APIsTimeTaken.StartQuiz), "Expected log to contain start quiz time")
	assert.Contains(t, logString, fmt.Sprintf("Submit Quiz: %d ms", result.APIsTimeTaken.SubmitQuiz), "Expected log to contain submit quiz time")
	assert.Contains(t, logString, fmt.Sprintf("Report API: %d ms", result.APIsTimeTaken.ReportAPI), "Expected log to contain report API time")
	assert.Contains(t, logString, fmt.Sprintf("Email API: %d ms", result.APIsTimeTaken.EmailAPI), "Expected log to contain email API time")
}

func Test_app_results_GetResultLog_WhenNoAPIsTimeTaken(t *testing.T) {
	startTime := time.Now().UnixMilli()
	endTime := startTime + 1000
	result := &Session{
		ID:        "",
		UserID:    "testuser",
		Email:     "test@example.com",
		Topic:     "test-topic",
		StartTime: startTime,
		EndTime:   endTime,
		Question: []quizapi.Question{
			{
				ID:       "q1",
				Question: "What is the capital of France?",
				Options:  []string{"Paris", "London", "Berlin", "Madrid"},
			},
		},
		Answers: []quizapi.Answer{
			{
				QuestionID: "q1",
				Answer:     "Paris",
			},
		},
		Status: STATUS_COMPLETED,
		Error:  errors.New("Test error"),
		Score:  100,
		Report: "./tmp/reports/a.pdf",
	}

	logString := getResultLog(result)
	require.NotEmpty(t, logString, "Expected getResultLog() to return a non-empty string")
	assert.Contains(t, logString, "Error while starting", "Expected log to contain Error indicator")
	assert.Contains(t, logString, result.Email, "Expected log to contain email")
	assert.Contains(t, logString, result.UserID, "Expected log to contain user ID")
	assert.Contains(t, logString, fmt.Sprintf("%d", result.Score), "Expected log to contain score")
	assert.Contains(t, logString, string(result.Status), "Expected log to contain status")
	assert.Contains(t, logString, time.UnixMilli(result.StartTime).Format(time.RFC3339), "Expected log to contain start time")
	assert.Contains(t, logString, time.UnixMilli(result.EndTime).Format(time.RFC3339), "Expected log to contain end time")
	assert.Contains(t, logString, fmt.Sprintf("%d ms", result.EndTime-result.StartTime), "Expected log to contain time taken")
	assert.Contains(t, logString, fmt.Sprintf("%v", result.Answers), "Expected log to contain answers")
	assert.Contains(t, logString, result.Report, "Expected log to contain report path")
	assert.Contains(t, logString, "Error: Test error", "Expected log to contain error message")
	assert.Contains(t, logString, "APIs Time Taken: Not available", "Expected log to indicate APIs time taken is not available")
}

func Test_app_results_GetSummaryLog(t *testing.T) {
	timetaken := []int64{1000, 2000, 1500, 3000}
	numOfUsers := 4
	expectedAvgTime := int64(1875)

	summaryLog := getSummaryLog(timetaken, numOfUsers)

	require.NotEmpty(t, summaryLog, "Expected getSummaryLog() to return a non-empty string")
	assert.Contains(t, summaryLog, fmt.Sprintf("%d", numOfUsers), "Expected log to contain number of users")
	assert.Contains(t, summaryLog, fmt.Sprintf("%d", expectedAvgTime), "Expected log to contain average time taken")
}

func Test_app_results_ListenForResults(t *testing.T) {
	tmpDirPath = "./test" // change path for the test environment
	expectedFilePath := fmt.Sprintf("%s/logs.txt", tmpDirPath)

	ssid := "1234"
	email := "test@example.com"

	defer func() {
		if err := os.RemoveAll(tmpDirPath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Error cleaning up the test files: %s", tmpDirPath)
		}
	}()

	// remove results files if it exists
	if err := os.RemoveAll(tmpDirPath); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Error while removing existing tmp dir for test")
	}

	app := NewTestApp()
	app.ResultListener.Add(1)
	go app.ListenForResults()

	session := &Session{ID: ssid, Email: email, Status: STATUS_FAILED, Error: fmt.Errorf("session error")}
	app.Results <- session
	close(app.Results)

	app.ResultListener.Wait()

	file, err := os.Open(expectedFilePath)
	require.NoError(t, err, "Expected to open results file without error")
	defer file.Close()

	content, err := os.ReadFile(expectedFilePath)
	require.NoError(t, err, "Expected to read results file without error")

	logString := string(content)
	require.NotEmpty(t, logString, "Expected results file to contain log data")
	assert.Contains(t, logString, "Session ID: "+ssid, "Expected log to contain session ID")
	assert.Contains(t, logString, "RESULTS", "Expected log to contain summary")
}
