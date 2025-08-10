package app

import (
	"fmt"
	"sync"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
)

type App struct {
	Config    *Config
	Wait      *sync.WaitGroup
	Mu        sync.Mutex
	QuizAPI   *quizapi.QuizAPI
	Results   chan *Session
	Errors    chan error
	Listeners sync.WaitGroup
}

func NewApp() *App {
	cfg := LoadConfig()
	quizApi := quizapi.NewQuizAPI(
		cfg.BaseURL,
		cfg.ReportServerBaseURL,
	)
	return &App{
		Config:  LoadConfig(),
		Wait:    &sync.WaitGroup{},
		Mu:      sync.Mutex{},
		QuizAPI: quizApi,
		Results: make(chan *Session, cfg.NumUsers),
		Errors:  make(chan error),
	}
}

func (app *App) Stop() {
	close(app.Results)
	close(app.Errors)
	// wait for the results and errors to be processed
	fmt.Println("Waiting for results and errors to be processed...")
	app.Listeners.Wait()
}
