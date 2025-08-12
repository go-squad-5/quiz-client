package app

import (
	"log"
	"os"
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
	InfoLogger     *log.Logger
	ErrorLogger    *log.Logger
	DebugLogger    *log.Logger
	ResultLogger   *log.Logger
}

func NewApp() *App {
	cfg := LoadConfig()
	quizApi := quizapi.NewQuizAPI(
		cfg.BaseURL,
		cfg.ReportServerBaseURL,
	)

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ltime)
	debugLog := log.New(os.Stdout, "DEBUG\t", log.Ltime)
	resultLog := log.New(os.Stdout, "RESULT\t", log.Ltime)

	return &App{
		Config:         LoadConfig(),
		Wait:           &sync.WaitGroup{},
		QuizAPI:        quizApi,
		Results:        make(chan *Session, cfg.NumUsers),
		Errors:         make(chan error),
		ResultListener: &sync.WaitGroup{},
		ErrorListener:  &sync.WaitGroup{},
		InfoLogger:     infoLog,
		ErrorLogger:    errorLog,
		DebugLogger:    debugLog,
		ResultLogger:   resultLog,
	}
}

func (app *App) Stop() {
	// wait for the results and errors to be processed
	app.InfoLogger.Println("Waiting for results and errors to be processed...")
	close(app.Errors)
	app.ErrorListener.Wait()
	close(app.Results)
	app.ResultListener.Wait()
}
