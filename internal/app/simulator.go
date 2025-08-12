package app

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
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
		app.InfoLogger.Println("GO ROUTINE started for user simulation: ", email, "on topic:", topic)
		// go app.FakeUserAction(email, topic)
		go app.SimulateUser(email, topic)
	}
}

func (app *App) SimulateUser(email, topic string) {
	defer func() {
		if r := recover(); r != nil {
			app.ErrorLogger.Println("Recovered from panic in user", email, ":", r)
		}
		app.InfoLogger.Println("GO ROUTINE FINISHED for user simulation:", email, "on topic:", topic)
		app.Wait.Done()
	}()

	app.InfoLogger.Printf("Simulating user action for email: %s, topic: %s\n", email, topic)

	startTime := time.Now()

	// Create a session
	app.InfoLogger.Println("Sending Request to create session for email:", email, "on topic:", topic)
	createStart := time.Now()
	ssid, err := app.QuizAPI.CreateSession(email, topic)
	createEnd := time.Now()
	app.InfoLogger.Printf("Session created for email: %s, topic: %s, session ID: %s\n", email, topic, ssid)
	if err != nil {
		app.ErrorLogger.Printf("Error creating session for email: %s, topic: %s, error: %v\n", email, topic, err)
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
	app.InfoLogger.Println("Sending Request to start quiz for session ID:", ssid, "on topic:", topic)
	startQuizStart := time.Now()
	questions, err := app.QuizAPI.StartQuiz(ssid, topic)
	startQuizEnd := time.Now()
	app.InfoLogger.Printf("Got questions for session ID: %s, topic: %s, questions: %d\n", ssid, topic, len(questions))
	aPIsTimeTaken.StartQuiz = startQuizEnd.UnixMilli() - startQuizStart.UnixMilli()
	if err != nil {
		app.ErrorLogger.Printf("Error starting quiz for session ID: %s, topic: %s, error: %v\n", ssid, topic, err)
		session.Error = err
		session.Status = STATUS_FAILED
		session.EndTime = time.Now().UnixMilli()
		app.Errors <- &SessionError{
			Session: session,
		}
		return
	}
	app.InfoLogger.Printf("Quiz started for session ID: %s, topic: %s, questions: %d\n", ssid, topic, len(questions))

	session.Question = questions

	// mark random answers to questions
	for _, question := range questions {
		numOptions := len(question.Options)
		if numOptions == 0 {
			app.ErrorLogger.Println("No options available for question ID:", question.ID)
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
	app.InfoLogger.Println("Sending Request to submit quiz for session ID:", ssid)
	submitStart := time.Now()
	score, err := app.QuizAPI.SubmitQuiz(ssid, session.Answers)
	submitEnd := time.Now()
	aPIsTimeTaken.SubmitQuiz = submitEnd.UnixMilli() - submitStart.UnixMilli()
	app.InfoLogger.Printf("Quiz submitted for session ID: %s, score: %d\n", ssid, score)
	if err != nil {
		app.ErrorLogger.Printf("Error submitting quiz for session ID: %s, error: %v\n", ssid, err)
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

	wg := &sync.WaitGroup{}

	wg.Add(1)
	app.InfoLogger.Println("GO ROUTINE Started to get report for session ID:", ssid)
	go func() {
		defer wg.Done()
		defer app.InfoLogger.Println("GO ROUTINE FINISHED for getting report for session ID:", ssid)
		// Get the report for the session
		app.InfoLogger.Println("Sending Request to get report for session ID:", ssid)
		reportStart := time.Now()
		report, err := app.QuizAPI.GetReport(ssid)
		reportEnd := time.Now()
		aPIsTimeTaken.ReportAPI = reportEnd.UnixMilli() - reportStart.UnixMilli()
		app.InfoLogger.Printf("Report received for session ID: %s, report: %+v\n", ssid, report)
		if err != nil {
			app.ErrorLogger.Printf("Error getting report for session ID: %s, error: %v\n", ssid, err)
			session.Error = err
			session.Status = STATUS_FAILED
			session.EndTime = time.Now().UnixMilli()
			app.Errors <- &SessionError{
				Session: session,
			}
			return
		}
		session.Report = report
	}()

	wg.Add(1)
	app.InfoLogger.Println("GO ROUTINE Started to get email report for session ID:", ssid)
	go func() {
		defer wg.Done()
		defer app.InfoLogger.Println("GO ROUTINE FINISHED for getting email report for session ID:", ssid)
		// Do email request
		app.InfoLogger.Println("Sending Request to get email report for session ID:", ssid)
		emailStart := time.Now()
		_, err = app.QuizAPI.GetEmailReport(ssid)
		emailEnd := time.Now()
    app.InfoLogger.Printf("Email Request Successful for session ID: %s\n", ssid)
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
	}()

	wg.Wait()

	endTime := time.Now()
	session.EndTime = endTime.UnixMilli()
	session.Status = STATUS_COMPLETED

	app.InfoLogger.Printf("Session completed for email: %s, topic: %s, session ID: %s, score: %d\n Sending session data to the results channel to log.", email, topic, ssid, score)
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
