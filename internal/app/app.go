package app

import (
	"fmt"
	"sync"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
)

type App struct {
	Config         *Config
	Wait           *sync.WaitGroup
	QuizAPI        *quizapi.QuizAPI
	Results        chan *Session
	Errors         chan error
	ResultListener *sync.WaitGroup
	ErrorListener  *sync.WaitGroup
}

func NewApp() *App {
	cfg := LoadConfig()
	quizApi := quizapi.NewQuizAPI(
		cfg.BaseURL,
		cfg.ReportServerBaseURL,
	)
	return &App{
		Config:         LoadConfig(),
		Wait:           &sync.WaitGroup{},
		QuizAPI:        quizApi,
		Results:        make(chan *Session, cfg.NumUsers),
		Errors:         make(chan error),
		ResultListener: &sync.WaitGroup{},
		ErrorListener:  &sync.WaitGroup{},
	}
}

func (app *App) Stop() {
	// wait for the results and errors to be processed
	fmt.Println("Waiting for results and errors to be processed...")
	close(app.Errors)
	app.ErrorListener.Wait()
	close(app.Results)
	app.ResultListener.Wait()
}
