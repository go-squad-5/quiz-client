package app

import "github.com/go-squad-5/quiz-load-test/internal/quizapi"

type STATUS string

const (
	STATUS_CREATED   STATUS = "created"
	STATUS_STARTED   STATUS = "started"
	STATUS_COMPLETED STATUS = "completed"
	STATUS_FAILED    STATUS = "failed"
)

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

type APIsTimeTaken struct {
	SessionCreation int64
	StartQuiz       int64
	SubmitQuiz      int64
	ReportAPI       int64
	EmailAPI        int64
}
