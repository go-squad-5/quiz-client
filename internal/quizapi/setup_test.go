package quizapi

import (
	"net/http"
	"time"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestQuizAPI(baseUrl, reportServerBaseUrl string, transportFunc RoundTripFunc) *QuizAPI {
	return &QuizAPI{
		client: &http.Client{
			Timeout:   60 * time.Second,
			Transport: transportFunc,
		},
		endpoints: endpoints{
			createSession:  baseUrl + "/session/create",
			startQuiz:      baseUrl + "/quiz/start",
			submitQuiz:     baseUrl + "/quiz/submit",
			getReport:      reportServerBaseUrl + "/sessions/%s/report",
			getEmailReport: reportServerBaseUrl + "/sessions/%s/email-report",
		},
	}
}

// func TestMain(m *testing.M) {
// 	// test environment
// 	return
// }
