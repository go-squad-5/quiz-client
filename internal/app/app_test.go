package app

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_app_NewApp(t *testing.T) {
	app := NewApp()
	require.NotNil(t, app, "NewApp should return a non-nil App instance")
	require.NotNil(t, app.Config, "App Config should not be nil")
	require.NotNil(t, app.Wait, "App WaitGroup should not be nil")
	require.NotNil(t, app.QuizAPI, "App QuizAPI should not be nil")
	require.NotNil(t, app.Results, "App Results channel should not be nil")
	require.NotNil(t, app.Errors, "App Errors channel should not be nil")
	require.NotNil(t, app.ResultListener, "App ResultListener WaitGroup should not be nil")
	require.NotNil(t, app.ErrorListener, "App ErrorListener WaitGroup should not be nil")
	require.NotNil(t, app.InfoLogger, "App InfoLogger should not be nil")
	require.NotNil(t, app.ErrorLogger, "App ErrorLogger should not be nil")
	require.NotNil(t, app.DebugLogger, "App DebugLogger should not be nil")
	require.NotNil(t, app.ResultLogger, "App ResultLogger should not be nil")
}

func Test_app_Stop(t *testing.T) {
	app := NewApp()
	require.NotPanics(t, func() {
		app.Stop()
	}, "Stop should not panic")
}

func Test_app_InfoLogger(t *testing.T) {
	app := NewTestApp()

	var buf bytes.Buffer
	app.InfoLogger.SetOutput(&buf)

	app.InfoLogger.Print("Hello World")
	app.InfoLogger.SetOutput(os.Stdout)
	output := buf.String()

	assert.Containsf(t, output, "INFO\t", "Expected 'INFO' in the log")
	assert.Containsf(t, output, "Hello World", "Expected 'Hello World' in the log")
}

func Test_app_ErrorLogger(t *testing.T) {
	app := NewTestApp()

	var buf bytes.Buffer
	app.ErrorLogger.SetOutput(&buf)

	app.ErrorLogger.Print("Hello World")
	app.ErrorLogger.SetOutput(os.Stdout)
	output := buf.String()

	assert.Containsf(t, output, "ERROR\t", "Expected 'INFO' in the log")
	assert.Containsf(t, output, "Hello World", "Expected 'Hello World' in the log")
}

func Test_app_DebugLogger(t *testing.T) {
	app := NewTestApp()

	var buf bytes.Buffer
	app.DebugLogger.SetOutput(&buf)

	app.DebugLogger.Print("Hello World")
	app.DebugLogger.SetOutput(os.Stdout)
	output := buf.String()

	assert.Containsf(t, output, "DEBUG\t", "Expected 'INFO' in the log")
	assert.Containsf(t, output, "Hello World", "Expected 'Hello World' in the log")
}

func Test_app_ResultLogger(t *testing.T) {
	app := NewTestApp()

	var buf bytes.Buffer
	app.ResultLogger.SetOutput(&buf)

	app.ResultLogger.Print("Hello World")
	app.ResultLogger.SetOutput(os.Stdout)
	output := buf.String()

	assert.Containsf(t, output, "RESULT\t", "Expected 'INFO' in the log")
	assert.Containsf(t, output, "Hello World", "Expected 'Hello World' in the log")
}
