package app

import (
	"testing"

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
