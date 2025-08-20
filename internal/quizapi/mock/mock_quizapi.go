package mock

import (
	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
	"github.com/stretchr/testify/mock"
)

type MockQuizAPI struct {
	mock.Mock
}

func (m *MockQuizAPI) CreateSession(email, topic string) (string, error) {
	args := m.Called(email, topic)
	return args.String(0), args.Error(1)
}

func (m *MockQuizAPI) StartQuiz(sessionId, topic string) ([]quizapi.Question, error) {
	args := m.Called(sessionId, topic)
	return args.Get(0).([]quizapi.Question), args.Error(1)
}

func (m *MockQuizAPI) SubmitQuiz(sessionId string, answers []quizapi.Answer) (int, error) {
	args := m.Called(sessionId, answers)
	return args.Int(0), args.Error(1)
}

func (m *MockQuizAPI) GetReport(sessionId string) (string, error) {
	args := m.Called(sessionId)
	return args.String(0), args.Error(1)
}

func (m *MockQuizAPI) GetEmailReport(sessionId string) (string, error) {
	args := m.Called(sessionId)
	return args.String(0), args.Error(1)
}
