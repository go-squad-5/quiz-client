package app

import (
	"log"
	"os"
	"sync"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi/mock"
)

func NewTestApp() *App {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ltime)
	debugLog := log.New(os.Stdout, "DEBUG\t", log.Ltime)
	resultLog := log.New(os.Stdout, "RESULT\t", log.Ltime)
	cfg := Config{
		BaseURL:             "http://localhost:8080",
		ReportServerBaseURL: "http://localhost:8070",
		NumUsers:            10,
	}
	quizApi := &mock.MockQuizAPI{}
	return &App{
		Config:         &cfg,
		Wait:           &sync.WaitGroup{},
		Results:        make(chan *Session, cfg.NumUsers),
		Errors:         make(chan error),
		ResultListener: &sync.WaitGroup{},
		ErrorListener:  &sync.WaitGroup{},
		InfoLogger:     infoLog,
		ErrorLogger:    errorLog,
		DebugLogger:    debugLog,
		ResultLogger:   resultLog,
		QuizAPI:        quizApi,
	}
}
