package app

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
)

func (app *App) StartSimulation() {
	for i := range app.Config.NumUsers {
		numEmails := len(EMAILS)
		numTopics := len(TOPICS)
		email := EMAILS[i%numEmails]
		topic := TOPICS[i%numTopics]
		app.Wait.Add(1)
		go app.FakeUserAction(email, topic)
		// go app.SimulateUser(email, topic)
	}
}

func (app *App) SimulateUser(email, topic string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in user", email, ":", r)
		}
		app.Wait.Done()
	}()

	log.Printf("Simulating user action for email: %s, topic: %s\n", email, topic)

	startTime := time.Now()

	// Create a session
	createStart := time.Now()
	ssid, err := app.QuizAPI.CreateSession(email, topic)
	createEnd := time.Now()
	if err != nil {
		app.Errors <- &StartSessionError{
			Email: email,
			Topic: topic,
			err:   err,
		}
		return
	}

	aPIsTimeTaken := &APIsTimeTaken{
		SessionCreation: createEnd.UnixMilli() - createStart.UnixMilli(),
		StartQuiz:       0,
		SubmitQuiz:      0,
		ReportAPI:       0,
		EmailAPI:        0,
	}

	session := &Session{
		ID:            ssid,
		UserID:        "user_" + email,
		Email:         email,
		Topic:         topic,
		Status:        STATUS_STARTED,
		StartTime:     startTime.UnixMilli(),
		EndTime:       0,
		Answers:       []quizapi.Answer{},
		Score:         0,
		CreatedAt:     startTime.Unix(),
		APIsTimeTaken: aPIsTimeTaken,
	}

	// Start the quiz , get questions
	startQuizStart := time.Now()
	questions, err := app.QuizAPI.StartQuiz(ssid)
	startQuizEnd := time.Now()
	aPIsTimeTaken.StartQuiz = startQuizEnd.UnixMilli() - startQuizStart.UnixMilli()
	if err != nil {
		session.Error = err
		session.Status = STATUS_FAILED
		session.EndTime = time.Now().UnixMilli()
		app.Errors <- &StartSessionError{
			Email: email,
			Topic: topic,
			err:   err,
		}
		return
	}

	session.Question = questions

	// mark random answers to questions
	for _, question := range questions {
		numOptions := len(question.Options)
		if numOptions == 0 {
			session.Error = fmt.Errorf("no options available for question ID: %s", question.ID)
			session.Status = STATUS_FAILED
			session.EndTime = time.Now().UnixMilli()
			app.Errors <- &SessionError{
				Session: session,
			}
			return
		}
		index := rand.Intn(numOptions) // random index
		session.Answers = append(session.Answers, quizapi.Answer{
			QuestionID: question.ID,
			Answer:     question.Options[index],
		})
	}

	// Submit the quiz with answers
	submitStart := time.Now()
	score, err := app.QuizAPI.SubmitQuiz(ssid, session.Answers)
	submitEnd := time.Now()
	aPIsTimeTaken.SubmitQuiz = submitEnd.UnixMilli() - submitStart.UnixMilli()
	if err != nil {
		session.Status = STATUS_FAILED
		session.Error = err
		session.EndTime = time.Now().UnixMilli()
		app.Errors <- &SessionError{
			Session: session,
		}
		return
	}

	// update session with score and end time
	session.Score = score

	// Get the report for the session
	reportStart := time.Now()
	report, err := app.QuizAPI.GetReport(ssid)
	reportEnd := time.Now()
	aPIsTimeTaken.ReportAPI = reportEnd.UnixMilli() - reportStart.UnixMilli()
	if err != nil {
		session.Error = err
		session.Status = STATUS_FAILED
		session.EndTime = time.Now().UnixMilli()
		app.Errors <- &SessionError{
			Session: session,
		}
		return
	}
	session.Report = report

	// Do email request
	emailStart := time.Now()
	_, err = app.QuizAPI.GetEmailReport(ssid)
	emailEnd := time.Now()
	session.APIsTimeTaken.EmailAPI = emailEnd.UnixMilli() - emailStart.UnixMilli()
	if err != nil {
		session.Error = err
		session.Status = STATUS_FAILED
		session.EndTime = time.Now().UnixMilli()
		app.Errors <- &SessionError{
			Session: session,
		}
		return
	}

	endTime := time.Now()
	session.EndTime = endTime.UnixMilli()
	session.Status = STATUS_COMPLETED

	app.Results <- session
}

func (app *App) FakeUserAction(email, topic string) {
	defer app.Wait.Done()
	log.Printf("Simulating user action for email: %s, topic: %s\n", email, topic)
	startTime := time.Now()
	time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
	endTime := time.Now()

	app.Results <- &Session{
		ID:        "somereandomstring",
		UserID:    "user_" + email,
		Email:     email,
		Topic:     topic,
		Status:    STATUS_COMPLETED,
		StartTime: startTime.UnixMilli(),
		EndTime:   endTime.UnixMilli(),
	}
}
