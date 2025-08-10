package app

import (
	"sync"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
)

type App struct {
	Config   *Config
	Wait     *sync.WaitGroup
  Mu       sync.Mutex
	QuizAPI  *quizapi.QuizAPI
  Sessions map[string]*Session
}

func NewApp() *App {
	cfg := LoadConfig()
	quizApi := quizapi.NewQuizAPI(
		cfg.BaseURL,
		cfg.ReportServerBaseURL,
	)
	return &App{
		Config:   LoadConfig(),
		Wait:     &sync.WaitGroup{},
		Sessions: make(map[string]*Session),
		Mu:       sync.Mutex{},
		QuizAPI:  quizApi,
	}
}
