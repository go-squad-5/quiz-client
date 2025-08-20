package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_app_errors_SessionError_Error(t *testing.T) {
	session := &Session{ID: "12345"}
	err := &SessionError{Session: session}
	expected := "Session error for session ID: 12345 "
	assert.Equalf(t, expected, err.Error(), "Expected error message '%s', but got '%s'", expected, err.Error())
}

func Test_app_errors_StartSessionError_Error(t *testing.T) {
	email := "test@example.com"
	topic := "test-topic"
	err := &StartSessionError{
		Email: email,
		Topic: topic,
		err:   nil,
	}
	expected := fmt.Sprint("Failed to start session for email: ", email,
		" on topic: ", topic, " \n Error: ", err.err, "\n")
	assert.Equalf(t, expected, err.Error(), "Expected error message '%s', but got '%s'", expected, err.Error())
}

func Test_app_errors_ListenForErrors_WhenStartSessionError(t *testing.T) {
	email := "test@example.com"

	app := NewTestApp()
	app.ErrorListener.Add(1)
	go app.ListenForErrors()

	app.Errors <- &StartSessionError{Email: email, Topic: "test-topic", err: fmt.Errorf("start session error")}
	close(app.Errors)

	// start session session error
	startResult := <-app.Results
	require.NotNil(t, startResult, "Expected start session error result to be non-nil")
	assert.Equal(t, "", startResult.ID, "Expected session ID to be empty for start session error")
	assert.Equal(t, email, startResult.Email, "Expected email to match")
	assert.Equal(t, "test-topic", startResult.Topic, "Expected topic to match")
	assert.Equal(t, STATUS_FAILED, startResult.Status, "Expected session status to be failed")
	require.NotNil(t, startResult.Error, "Expected error to be non-nil for start session error")

	app.ErrorListener.Wait()
}

func Test_app_errors_ListenForErrors_WhenSessionError(t *testing.T) {
	ssid := "12345"
	email := "test@example.com"

	app := NewTestApp()
	app.ErrorListener.Add(1)
	go app.ListenForErrors()

	session := &Session{ID: ssid, Email: email, Status: STATUS_FAILED, Error: fmt.Errorf("session error")}

	app.Errors <- &SessionError{Session: session}
	close(app.Errors)

	// session error
	sessionResult := <-app.Results
	require.NotNil(t, sessionResult, "Expected session error result to be non-nil")
	assert.Equal(t, ssid, sessionResult.ID, "Expected session ID to match")
	assert.Equal(t, STATUS_FAILED, sessionResult.Status, "Expected session status to be failed")
	require.NotNil(t, sessionResult.Error, "Expected error to be non-nil for session error")

	app.ErrorListener.Wait()
}

func Test_app_errors_ListenForErrors_WhenUnknownError(t *testing.T) {
	app := NewTestApp()
	app.ErrorListener.Add(1)
	go app.ListenForErrors()

	app.Errors <- fmt.Errorf("generic error")
	close(app.Errors)

	// generic error
	genericResult := <-app.Results
	require.NotNil(t, genericResult, "Expected generic error result to be non-nil")
	assert.Equal(t, "", genericResult.ID, "Expected session ID to be empty for generic error")
	assert.Equal(t, "", genericResult.Email, "Expected email to be empty for generic error")
	assert.Equal(t, STATUS_FAILED, genericResult.Status, "Expected session status to be failed for generic error")
	require.NotNil(t, genericResult.Error, "Expected error to be non-nil for generic error")

	app.ErrorListener.Wait()
}
