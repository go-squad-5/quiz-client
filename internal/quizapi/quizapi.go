package quizapi

import (
	"net/http"
	"time"
)

type IQuizAPI interface {
	CreateSession(email, topic string) (string, error)
	StartQuiz(sessionId, topic string) ([]Question, error)
	SubmitQuiz(sessionId string, answers []Answer) (int, error) // Score, error
	GetReport(sessionId string) (string, error)
	GetEmailReport(sessionId string) (string, error)
}

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
