package app

import (
	"time"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
)

type STATUS string

const (
	STATUS_CREATED   STATUS = "created"
	STATUS_STARTED   STATUS = "started"
	STATUS_COMPLETED STATUS = "completed"
	STATUS_FAILED    STATUS = "failed"
)

type APIsTimeTaken struct {
	SessionCreation int64
	StartQuiz       int64
	SubmitQuiz      int64
	ReportAPI       int64
	EmailAPI        int64
}

func NewAPIsTimeTaken() *APIsTimeTaken {
	return &APIsTimeTaken{}
}

func (a *APIsTimeTaken) SetSessionCreationTime(timetaken int64) {
	a.SessionCreation = timetaken
}

func (a *APIsTimeTaken) SetStartQuizTime(timetaken int64) {
	a.StartQuiz = timetaken
}

func (a *APIsTimeTaken) SetSubmitQuizTime(timetaken int64) {
	a.SubmitQuiz = timetaken
}

func (a *APIsTimeTaken) SetReportAPITime(timetaken int64) {
	a.ReportAPI = timetaken
}

func (a *APIsTimeTaken) SetEmailAPITime(timetaken int64) {
	a.EmailAPI = timetaken
}

type Session struct {
	ID            string
	Email         string
	Topic         string
	UserID        string
	StartTime     int64
	EndTime       int64
	Question      []quizapi.Question
	Answers       []quizapi.Answer
	Score         int
	Report        string
	Status        STATUS
	Error         error
	CreatedAt     int64
	APIsTimeTaken *APIsTimeTaken
}

func NewSession(email, topic string, aPIsTimeTaken *APIsTimeTaken) *Session {
	return &Session{
		ID:            "",
		UserID:        "user_" + email,
		Email:         email,
		Topic:         topic,
		Status:        STATUS_STARTED,
		StartTime:     time.Now().UnixMilli(),
		EndTime:       0,
		Answers:       []quizapi.Answer{},
		Score:         0,
		CreatedAt:     time.Now().UnixMilli(),
		APIsTimeTaken: aPIsTimeTaken,
	}
}

func (s *Session) SetSession(ssid string) {
	s.ID = ssid
}

func (s *Session) SetStatus(status STATUS) {
	s.Status = status
}

func (s *Session) SetStartTime(timestamp time.Time) {
	s.StartTime = timestamp.UnixMilli()
}

func (s *Session) SetEndTime(timestamp time.Time) {
	s.EndTime = timestamp.UnixMilli()
}

func (s *Session) SetAnswers(answers []quizapi.Answer) {
	s.Answers = answers
}

func (s *Session) SetScore(score int) {
	s.Score = score
}

func (s *Session) SetReport(report string) {
	s.Report = report
}

func (s *Session) SetError(err error) {
	s.Error = err
}

func (s *Session) SetQuestions(questions []quizapi.Question) {
	s.Question = questions
}
