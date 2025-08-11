package quizapi

import (
	"net/http"
	"time"
)

type QuizAPI struct {
	client    *http.Client
	endpoints endpoints
}

type endpoints struct {
	createSession  string
	startQuiz      string
	submitQuiz     string
	getReport      string
	getEmailReport string
}

func NewQuizAPI(baseUrl, reportServerBaseUrl string) *QuizAPI {
	return &QuizAPI{
		client: &http.Client{
			Timeout: 60 * time.Second,
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
