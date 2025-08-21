package app

import (
	"errors"
	"testing"
	"time"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
	"github.com/go-squad-5/quiz-load-test/internal/quizapi/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_app_simulator_CallCreateSession_WhenSuccess(t *testing.T) {
	email := "test@example.com"
	topic := "math"
	expectedSsid := "12345"

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	call := mockApp.On("CreateSession", email, topic).
		After(time.Millisecond*100). // 100 ms
		Return(expectedSsid, nil)

	ssid, timeTaken, err := app.callCreateSession(email, topic)
	require.NoError(t, err, "Expected create session call to return no error")
	mockApp.AssertExpectations(t)
	assert.GreaterOrEqual(t, int64(timeTaken), int64(100), "Expected create session api call time to be at least 1 second")
	assert.Equal(t, expectedSsid, ssid, "Expected create session api call to return the correct session id")
	call.Unset()
}

func Test_app_simulator_CallCreateSession_WhenFailure(t *testing.T) {
	email := "test@example.com"
	topic := "math"
	expectedSsid := ""
	expectedError := errors.New("failed to create session")

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	call := mockApp.On("CreateSession", email, topic).
		Return(expectedSsid, expectedError).
		After(time.Millisecond * 100)

	// test if errors are sent correctly
	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		err := <-app.Errors
		require.NotNil(t, err, "Expected to get non-nil error from Errors channel for failure case")
		_, ok := err.(*StartSessionError)
		require.Truef(t, ok, "Expected error to be the start session error")
	}()

	ssid, timeTaken, err := app.callCreateSession(email, topic)

	mockApp.AssertExpectations(t)
	require.Error(t, err, "Expected an error when failure case")
	assert.GreaterOrEqual(t, int64(timeTaken), int64(100), "Expected create session api call time to be at least 1 second")
	assert.Equal(t, expectedSsid, ssid, "Expected create session api call to return the correct session id")
	call.Unset()
	app.Wait.Wait()
}

func Test_app_simulator_CallStartQuiz_Success(t *testing.T) {
	ssid := "12345"
	topic := "math"
	expectedSsid := "12345"
	expectedQuestions := []quizapi.Question{
		{
			ID:       "q1",
			Question: "question one",
			Options:  []string{"1", "2", "3", "4"},
		},
	}

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	mockApp.On("StartQuiz", ssid, topic).
		After(time.Millisecond*100). // 100 ms
		Return(expectedQuestions, nil)

	session := NewSession("mohit@mohit.com", topic, nil)
	questions, timeTaken, err := app.callStartQuiz(ssid, topic, session)
	require.NoError(t, err, "Expected start quiz call to return no error")
	mockApp.AssertExpectations(t)
	assert.GreaterOrEqual(t, int64(timeTaken), int64(100), "Expected start quiz api call time to be at least 1 second")
	assert.Equal(t, expectedSsid, ssid, "Expected create session api call to return the correct session id")
	assert.Equal(t, expectedQuestions, questions, "Expected create session api call to return the correct Questions")
}

func Test_app_simulator_CallStartQuiz_WhenFailure(t *testing.T) {
	ssid := "12345"
	topic := "math"
	expectedError := errors.New("failed to start session")

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	mockApp.On("StartQuiz", ssid, topic).
		After(time.Millisecond*100). // 100 ms
		Return([]quizapi.Question{}, expectedError)

	// test if errors are sent correctly
	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		err := <-app.Errors
		require.NotNil(t, err, "Expected to get non-nil error from Errors channel for failure case")
		_, ok := err.(*SessionError)
		require.Truef(t, ok, "Expected error to be the session error")
	}()

	session := NewSession("mohit@mohit.com", topic, nil)
	questions, timeTaken, err := app.callStartQuiz(ssid, topic, session)
	close(app.Errors)

	mockApp.AssertExpectations(t)
	require.Error(t, err, "Expected an error when failure case")
	assert.GreaterOrEqual(t, int64(timeTaken), int64(100), "Expected create session api call time to be at least 1 second")
	assert.Empty(t, questions, "Expected Empty/Nil Questions to be returned")
	app.Wait.Wait()
}

func Test_app_simulator_CallStartQuiz_WhenNilSession(t *testing.T) {
	app := NewTestApp()
	_, _, err := app.callStartQuiz("1234", "math", nil)
	require.Error(t, err, "Expected callStartQuiz to return error when passing nil session value")
}

func Test_app_simulator_CallMarkRandomAnswers_WhenNilSession(t *testing.T) {
	app := NewTestApp()
	err := app.markRandomAnswers(nil, nil)
	require.Error(t, err, "Expected callStartQuiz to return error when passing nil session value")
}

func Test_app_simulator_CallMarkRandomAnswers_WhenNilQuestions(t *testing.T) {
	app := NewTestApp()
	session := NewSession("mohit@mohit.com", "", nil)
	err := app.markRandomAnswers(nil, session)
	require.Error(t, err, "Expected callStartQuiz to return error when passing nil questions value")
}

func Test_app_simulator_CallMarkRandomAnswers_WhenNoOptions(t *testing.T) {
	app := NewTestApp()
	session := NewSession("mohit@mohit.com", "", nil)
	questions := []quizapi.Question{
		{
			ID:       "q1",
			Question: "question?",
			Options:  []string{},
		},
	}
	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		err := <-app.Errors
		require.Error(t, err, "Expected to get a Error in no options case")
		_, ok := err.(*SessionError)
		require.True(t, ok, "Expected error to be a session error type")
	}()
	err := app.markRandomAnswers(questions, session)
	require.Error(t, err, "Expected callStartQuiz to return error when passing empty questions options value")
	app.Wait.Wait()
}

func Test_app_simulator_CallMarkRandomAnswers_WhenSuccess(t *testing.T) {
	app := NewTestApp()
	session := NewSession("mohit@mohit.com", "", nil)
	questions := []quizapi.Question{
		{
			ID:       "q1",
			Question: "question?",
			Options:  []string{"1", "2", "3"},
		},
	}
	err := app.markRandomAnswers(questions, session)

	require.NoError(t, err, "Expected callStartQuiz to return error when passing valid inputs")
	assert.NotEmpty(t, session.Answers, "Expected a non-empty answers in the session")
	assert.NotEmpty(t, session.Answers[0].Answer, "Expected a non-empty answers in the session, for each question")
	assert.NotEmpty(t, session.Answers[0].QuestionID, "Expected a non-empty questionID in the session, for each question")
}

func Test_app_simulator_CallSubmitQuiz_WhenNilSession(t *testing.T) {
	app := NewTestApp()
	_, _, err := app.callSubmitQuiz("1234", nil)
	require.Error(t, err, "Expected callSubmitQuiz to return error when passing nil session value")
}

func Test_app_simulator_CallSubmitQuiz_WhenSuccess(t *testing.T) {
	ssid := "12345"
	session := NewSession("test@example.com", "math", nil)
	answers := []quizapi.Answer{
		{
			QuestionID: "q1",
			Answer:     "1",
		},
	}
	session.SetAnswers(answers)
	expectedScore := 10
	expectedTimeTaken := int64(100)

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	call := mockApp.On("SubmitQuiz", ssid, answers).
		After(time.Millisecond*100). // 100 ms
		Return(expectedScore, nil)

	score, timeTaken, err := app.callSubmitQuiz(ssid, session)

	require.NoError(t, err, "Expected callSubmitQuiz to return no error")
	mockApp.AssertExpectations(t)
	assert.GreaterOrEqual(t, int64(timeTaken), expectedTimeTaken, "Expected callSubmitQuiz api call time to be at least 100 ms")
	assert.Equal(t, expectedScore, score, "Expected callSubmitQuiz api call to return the correct score")
	call.Unset()
}

func Test_app_simulator_CallSubmitQuiz_WhenError(t *testing.T) {
	ssid := "12345"
	session := NewSession("test@example.com", "math", nil)
	answers := []quizapi.Answer{
		{
			QuestionID: "q1",
			Answer:     "1",
		},
	}
	session.SetAnswers(answers)
	expectedScore := 0
	expectedTimeTaken := int64(100)
	expectedError := errors.New("failed to submit quiz")

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	call := mockApp.On("SubmitQuiz", ssid, answers).
		After(time.Millisecond*100). // 100 ms
		Return(expectedScore, expectedError)

	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		err := <-app.Errors
		require.NotNil(t, err, "Expected to get non-nil error from Errors channel for failure case")
		_, ok := err.(*SessionError)
		require.Truef(t, ok, "Expected error to be the session error")
	}()

	score, timeTaken, err := app.callSubmitQuiz(ssid, session)

	require.Error(t, err, "Expected callSubmitQuiz to return an error")
	mockApp.AssertExpectations(t)
	assert.GreaterOrEqual(t, int64(timeTaken), expectedTimeTaken, "Expected callSubmitQuiz api call time to be at least 100 ms")
	assert.Equal(t, expectedScore, score, "Expected callSubmitQuiz api call to return the correct score")
	app.Wait.Wait()
	call.Unset()
}

func Test_app_simulator_CallGetReport_WhenSuccess(t *testing.T) {
	session := NewSession("test@example.com", "math", nil)
	expectedReport := "This is a test report"
	expectedTimeTaken := int64(100)
	app := NewTestApp()
	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}
	call := mockApp.On("GetReport", session.ID).
		After(time.Millisecond*100).
		Return(expectedReport, nil)
	report, timeTaken, err := app.callGetReport(session)
	require.NoError(t, err, "Expected callGetReport to return no error")
	mockApp.AssertExpectations(t)
	assert.GreaterOrEqual(t, int64(timeTaken), expectedTimeTaken, "Expected callGetReport api call time to be at least 100 ms")
	assert.Equal(t, expectedReport, report, "Expected callGetReport api call to return the correct report")
	call.Unset()
}

func Test_app_simulator_CallGetReport_WhenError(t *testing.T) {
	session := NewSession("test@example.com", "math", nil)
	expectedReport := ""
	expectedTimeTaken := int64(100)
	expectedError := errors.New("failed to get report")

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	call := mockApp.On("GetReport", session.ID).
		After(time.Millisecond*100).
		Return(expectedReport, expectedError)

	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		err := <-app.Errors
		require.NotNil(t, err, "Expected to get non-nil error from Errors channel for failure case")
		_, ok := err.(*SessionError)
		require.Truef(t, ok, "Expected error to be the session error")
	}()

	report, timeTaken, err := app.callGetReport(session)

	require.Error(t, err, "Expected callGetReport to return an error")
	mockApp.AssertExpectations(t)
	assert.GreaterOrEqual(t, int64(timeTaken), expectedTimeTaken, "Expected callGetReport api call time to be at least 100 ms")
	assert.Equal(t, expectedReport, report, "Expected callGetReport api call to return the correct report")

	app.Wait.Wait()
	call.Unset()
}

func Test_app_simulator_CallGetReport_WhenNilSession(t *testing.T) {
	app := NewTestApp()
	_, _, err := app.callGetReport(nil)
	require.Error(t, err, "Expected callGetReport to return error when passing nil session value")
}

func Test_app_simulator_CallGetEmail_WhenSuccess(t *testing.T) {
	session := NewSession("test@example.com", "math", nil)
	expectedTimeTaken := int64(100)

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	call := mockApp.On("GetEmailReport", session.ID).
		After(time.Millisecond*100).
		Return("", nil)

	timeTaken, err := app.callGetEmail(session)

	require.NoError(t, err, "Expected callGetEmail to return no error")
	mockApp.AssertExpectations(t)
	assert.GreaterOrEqual(t, int64(timeTaken), expectedTimeTaken, "Expected callGetEmail api call time to be at least 100 ms")
	call.Unset()
}

func Test_app_simulator_CallGetEmail_WhenError(t *testing.T) {
	session := NewSession("test@example.com", "math", nil)
	expectedTimeTaken := int64(100)
	expectedError := errors.New("failed to get email report")

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)

	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	call := mockApp.On("GetEmailReport", session.ID).
		After(time.Millisecond*100).
		Return("", expectedError)

	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		err := <-app.Errors
		require.NotNil(t, err, "Expected to get non-nil error from Errors channel for failure case")
		_, ok := err.(*SessionError)
		require.Truef(t, ok, "Expected error to be the session error")
	}()

	timeTaken, err := app.callGetEmail(session)
	require.Error(t, err, "Expected callGetEmail to return an error")
	mockApp.AssertExpectations(t)
	assert.GreaterOrEqual(t, int64(timeTaken), expectedTimeTaken, "Expected callGetEmail api call time to be at least 100 ms")

	app.Wait.Wait()
	call.Unset()
}

func Test_app_simulator_CallGetEmail_WhenNilSession(t *testing.T) {
	app := NewTestApp()
	_, err := app.callGetEmail(nil)
	require.Error(t, err, "Expected callGetEmail to return error when passing nil session value")
}

func Test_app_simulator_CallReportAndEmailAPIs(t *testing.T) {
	apisTimeTaken := &APIsTimeTaken{}
	session := NewSession("test@example.com", "math", apisTimeTaken)

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	mockApp.On("GetReport", session.ID).
		After(time.Millisecond*100).
		Return("This is a test report", nil)
	mockApp.On("GetEmailReport", session.ID).
		After(time.Millisecond*100).
		Return("", nil)

	app.callReportAndEmailAPIs(session)

	require.NotEmpty(t, session.Report, "Expected session to have a non-nil report after calling callReportAndEmailAPIs")
	require.NotEmpty(t, session.APIsTimeTaken.ReportAPI, "Expected session to have a non-nil ReportAPI time after calling callReportAndEmailAPIs")
	require.NotEmpty(t, session.APIsTimeTaken.EmailAPI, "Expected session to have a non-nil EmailAPI time after calling callReportAndEmailAPIs")
}

func Test_app_simulator_callReportAndEmailAPIs_WhenNilSession(t *testing.T) {
	app := NewTestApp()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected callReportAndEmailAPIs to panic when passing nil session")
		}
	}()
	app.callReportAndEmailAPIs(nil)
}

func Test_app_simulator_SimulateUser_WhenSuccess(t *testing.T) {
	email := "test@example.com"
	topic := "math"

	app := NewTestApp()

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	// write all expectations for the mock app
	mockApp.On("CreateSession", email, topic).
		After(time.Millisecond*100). // 100 ms
		Return("12345", nil)
	mockApp.On("StartQuiz", "12345", topic).
		After(time.Millisecond*100). // 100 ms
		Return([]quizapi.Question{
			{
				ID:       "q1",
				Question: "What is 2 + 2?",
				Options:  []string{"3"},
			},
		}, nil)
	mockApp.On("SubmitQuiz", "12345", []quizapi.Answer{
		{
			QuestionID: "q1",
			Answer:     "3",
		},
	}).After(time.Millisecond*100). // 100 ms
					Return(10, nil)
	mockApp.On("GetReport", "12345").
		After(time.Millisecond*100). // 100 ms
		Return("This is a test report", nil)
	mockApp.On("GetEmailReport", "12345").
		After(time.Millisecond*100). // 100 ms
		Return("", nil)

	// close errors channel to avoid deadlock
	close(app.Errors)
	// ensure no one sends to errors channel
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, but got: %v", r)
		}
	}()

	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		result := <-app.Results
		require.NotNil(t, result, "Expected to get a non-nil result from Results channel")
		assert.Equal(t, "12345", result.ID, "Expected result ID to match the session ID")
		assert.Equal(t, email, result.Email, "Expected result Email to match the session Email")
		assert.Equal(t, topic, result.Topic, "Expected result Topic to match the session Topic")
		assert.NotEmpty(t, result.Question, "Expected result Question to be non-empty")
		assert.NotEmpty(t, result.Answers, "Expected result Answers to be non-empty")
		assert.Equal(t, 10, result.Score, "Expected result Score to match the expected score")
		assert.NotEmpty(t, result.Report, "Expected result Report to be non-empty")
		assert.NotEmpty(t, result.APIsTimeTaken.SessionCreation, "Expected result APIsTimeTaken SessionCreation to be non-empty")
		assert.NotEmpty(t, result.APIsTimeTaken.StartQuiz, "Expected result APIsTimeTaken StartQuiz to be non-empty")
		assert.NotEmpty(t, result.APIsTimeTaken.SubmitQuiz, "Expected result APIsTimeTaken SubmitQuiz to be non-empty")
		assert.NotEmpty(t, result.APIsTimeTaken.ReportAPI, "Expected result APIsTimeTaken ReportAPI to be non-empty")
		assert.NotEmpty(t, result.APIsTimeTaken.EmailAPI, "Expected result APIsTimeTaken EmailAPI to be non-empty")
		assert.NoError(t, result.Error, "Expected result Error to be nil")
	}()

	app.Wait.Add(1)
	app.SimulateUser(email, topic)

	app.Wait.Wait()
	mockApp.AssertExpectations(t)
}

func Test_app_simulator_StartSimulation(t *testing.T) {
	app := NewTestApp() // default 10 users

	mockApp, ok := app.QuizAPI.(*mock.MockQuizAPI)
	if !ok {
		t.Fatal("Error while getting the mock quizapi")
	}

	for i := range app.Config.NumUsers {
		email := EMAILS[i%len(EMAILS)]
		topic := TOPICS[i%len(TOPICS)]
		// write all expectations for the mock app
		mockApp.On("CreateSession", email, topic).
			After(time.Millisecond*100). // 100 ms
			Return("12345", nil)
		mockApp.On("StartQuiz", "12345", topic).
			After(time.Millisecond*100). // 100 ms
			Return([]quizapi.Question{
				{
					ID:       "q1",
					Question: "What is 2 + 2?",
					Options:  []string{"3"},
				},
			}, nil)
		mockApp.On("SubmitQuiz", "12345", []quizapi.Answer{
			{
				QuestionID: "q1",
				Answer:     "3",
			},
		}).After(time.Millisecond*100). // 100 ms
						Return(10, nil)
		mockApp.On("GetReport", "12345").
			After(time.Millisecond*100). // 100 ms
			Return("This is a test report", nil)
		mockApp.On("GetEmailReport", "12345").
			After(time.Millisecond*100). // 100 ms
			Return("", nil)
	}

	// close errors channel to avoid deadlock
	close(app.Errors)
	// ensure no one sends to errors channel
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, but got: %v", r)
		}
	}()

	done := make(chan int)
	go func() {
		count := 0
		for range app.Config.NumUsers {
			result := <-app.Results
			count++
			require.NotNil(t, result, "Expected to get a non-nil result from Results channel")
			assert.Equal(t, "12345", result.ID, "Expected result ID to match the session ID")
			assert.NotEmpty(t, result.Email, "Expected result Email to match the session Email")
			assert.NotEmpty(t, result.Topic, "Expected result Topic to match the session Topic")
			assert.NotEmpty(t, result.Question, "Expected result Question to be non-empty")
			assert.NotEmpty(t, result.Answers, "Expected result Answers to be non-empty")
			assert.Equal(t, 10, result.Score, "Expected result Score to match the expected score")
			assert.NotEmpty(t, result.Report, "Expected result Report to be non-empty")
			assert.NotEmpty(t, result.APIsTimeTaken.SessionCreation, "Expected result APIsTimeTaken SessionCreation to be non-empty")
			assert.NotEmpty(t, result.APIsTimeTaken.StartQuiz, "Expected result APIsTimeTaken StartQuiz to be non-empty")
			assert.NotEmpty(t, result.APIsTimeTaken.SubmitQuiz, "Expected result APIsTimeTaken SubmitQuiz to be non-empty")
			assert.NotEmpty(t, result.APIsTimeTaken.ReportAPI, "Expected result APIsTimeTaken ReportAPI to be non-empty")
			assert.NotEmpty(t, result.APIsTimeTaken.EmailAPI, "Expected result APIsTimeTaken EmailAPI to be non-empty")
			assert.NoError(t, result.Error, "Expected result Error to be nil")
		}
		done <- count
	}()

	app.StartSimulation()
	app.Wait.Wait()
	close(app.Results)
	if count := <-done; count != app.Config.NumUsers {
		t.Errorf("Expected to get %d results, but got %d", app.Config.NumUsers, done)
	}
}
