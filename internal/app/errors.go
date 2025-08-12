package app

import "fmt"

type SessionError struct {
	Session *Session
}

func (e *SessionError) Error() string {
	return fmt.Sprintf("Session error for session ID: %s ", e.Session.ID)
}

type StartSessionError struct {
	Email string
	Topic string
	err   error
}

func (e *StartSessionError) Error() string {
	return fmt.Sprint("Failed to start session for email: ", e.Email, " on topic: ", e.Topic, " \n Error: ", e.err, "\n")
}

func (app *App) ListenForErrors() {
	defer app.ErrorListener.Done()
	defer app.InfoLogger.Println("GO ROUTINE FINISHED for listening to errors")
	for err := range app.Errors {
		if sessionErr, ok := err.(*SessionError); ok {
			// Handle session-specific errors
			app.Results <- sessionErr.Session
		} else if startErr, ok := err.(*StartSessionError); ok {
			// Handle start session errors
			app.Results <- &Session{
				ID:     "",
				Email:  startErr.Email,
				Status: STATUS_FAILED,
				Error:  err,
			}
		} else {
			// Handle general errors
			app.Results <- &Session{
				ID:     "",
				Status: STATUS_FAILED,
				Error:  err,
			}
		}
	}
}
