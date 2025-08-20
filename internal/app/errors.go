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
	return fmt.Sprint("Failed to start session for email: ", e.Email,
		" on topic: ", e.Topic, " \n Error: ", e.err, "\n")
}

func (app *App) ListenForErrors() {
	defer app.ErrorListener.Done()
	defer app.InfoLogger.Println("GO ROUTINE FINISHED for listening to errors")
	for err := range app.Errors {
		switch e := err.(type) {
		case *StartSessionError:
			app.Results <- &Session{
				ID:     "",
				Email:  e.Email,
				Topic:  e.Topic,
				Status: STATUS_FAILED,
				Error:  err,
			}
		case *SessionError:
			app.Results <- e.Session
		default:
			app.Results <- &Session{
				ID:     "",
				Status: STATUS_FAILED,
				Error:  err,
			}
		}
	}
}
