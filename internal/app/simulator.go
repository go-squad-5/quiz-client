package app

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-squad-5/quiz-load-test/internal/quizapi"
)

func (app *App) StartSimulation() {
	for i := range app.Config.NumUsers {
		numEmails, numTopics := getNumberOfEmailsAndTopics()
		email := EMAILS[i%numEmails]
		topic := TOPICS[i%numTopics]
		app.Wait.Add(1)
		app.InfoLogger.Println("GO ROUTINE started for user simulation: ", email, "on topic:", topic)
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

	// create session struct
	aPIsTimeTaken := NewAPIsTimeTaken()
	session := NewSession(email, topic, aPIsTimeTaken)

	ssid, createTimeTaken, err := app.callCreateSession(email, topic)
	aPIsTimeTaken.SetSessionCreationTime(createTimeTaken)
	if err != nil {
		return
	}
	session.SetSession(ssid)

	questions, startQuizTimeTaken, err := app.callStartQuiz(ssid, topic, session)
	aPIsTimeTaken.SetStartQuizTime(startQuizTimeTaken)
	if err != nil {
		return
	}
	session.SetQuestions(questions)

	// mark random answers to questions
	if err := app.markRandomAnswers(questions, session); err != nil {
		return
	}

	score, submitTimeTaken, err := app.callSubmitQuiz(ssid, session)
	aPIsTimeTaken.SetSubmitQuizTime(submitTimeTaken)
	session.SetScore(score)

	// call report and email apis concurrently
	app.callReportAndEmailAPIs(session)

	// end session
	session.SetEndTime(time.Now())
	session.SetStatus(STATUS_COMPLETED)

	app.InfoLogger.Printf("Session completed for email: %s, topic: %s, session ID: %s, score: %d\n Sending session data to the results channel to log.", email, topic, ssid, score)
	app.Results <- session
}

func (app *App) callCreateSession(email, topic string) (string, int64, error) {
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
		return "", createEnd.UnixMilli() - createStart.UnixMilli(), err
	}
	return ssid, createEnd.UnixMilli() - createStart.UnixMilli(), nil
}

func (app *App) callStartQuiz(ssid, topic string, session *Session) ([]quizapi.Question, int64, error) {
	if session == nil {
		return nil, 0, fmt.Errorf("sesssion should be non-nil value")
	}
	app.InfoLogger.Println("Sending Request to start quiz for session ID:", ssid, "on topic:", topic)
	startQuizStart := time.Now()
	questions, err := app.QuizAPI.StartQuiz(ssid, topic)
	startQuizEnd := time.Now()
	app.InfoLogger.Printf("Got questions for session ID: %s, topic: %s, questions: %d\n", ssid, topic, len(questions))
	if err != nil {
		app.ErrorLogger.Printf("Error starting quiz for session ID: %s, topic: %s, error: %v\n", ssid, topic, err)
		session.SetError(err)
		session.SetStatus(STATUS_FAILED)
		session.SetEndTime(time.Now())
		app.Errors <- &SessionError{
			Session: session,
		}
		return nil, startQuizEnd.UnixMilli() - startQuizStart.UnixMilli(), err
	}
	app.InfoLogger.Printf("Quiz started for session ID: %s, topic: %s, questions: %d\n", ssid, topic, len(questions))
	return questions, startQuizEnd.UnixMilli() - startQuizStart.UnixMilli(), nil
}

func (app *App) markRandomAnswers(questions []quizapi.Question, session *Session) error {
	if session == nil {
		return fmt.Errorf("sesssion should be non-nil value")
	}
	if questions == nil {
		return fmt.Errorf("questions slice should be non-nil")
	}
	answers := make([]quizapi.Answer, 0, len(questions))
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
			return session.Error
		}
		index := rand.Intn(numOptions) // random index
		answers = append(answers, quizapi.Answer{
			QuestionID: question.ID,
			Answer:     question.Options[index],
		})
	}
	session.SetAnswers(answers)
	return nil
}

func (app *App) callSubmitQuiz(ssid string, session *Session) (int, int64, error) {
	if session == nil {
		return 0, 0, fmt.Errorf("sesssion should be non-nil value")
	}
	app.InfoLogger.Println("Sending Request to submit quiz for session ID:", ssid)
	submitStart := time.Now()
	score, err := app.QuizAPI.SubmitQuiz(ssid, session.Answers)
	submitEnd := time.Now()
	app.InfoLogger.Printf("Quiz submitted for session ID: %s, score: %d\n", ssid, score)
	if err != nil {
		app.ErrorLogger.Printf("Error submitting quiz for session ID: %s, error: %v\n", ssid, err)
		session.Status = STATUS_FAILED
		session.Error = err
		session.EndTime = time.Now().UnixMilli()
		app.Errors <- &SessionError{
			Session: session,
		}
		return 0, submitEnd.UnixMilli() - submitStart.UnixMilli(), err
	}
	return score, submitEnd.UnixMilli() - submitStart.UnixMilli(), nil
}

func (app *App) callReportAndEmailAPIs(session *Session) {
	if session == nil {
		panic("session should be non-nil value")
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	app.InfoLogger.Println("GO ROUTINE Started to get report for session ID:", session.ID)
	go func() {
		defer wg.Done()
		defer app.InfoLogger.Println("GO ROUTINE FINISHED for getting report for session ID:", session.ID)
		// Get the report for the session
		report, reportTimeTaken, _ := app.callGetReport(session)
		session.APIsTimeTaken.SetReportAPITime(reportTimeTaken)
		session.SetReport(report)
	}()

	wg.Add(1)
	app.InfoLogger.Println("GO ROUTINE Started to get email report for session ID:", session.ID)
	go func() {
		defer wg.Done()
		defer app.InfoLogger.Println("GO ROUTINE FINISHED for getting email report for session ID:", session.ID)
		// Do email request
		timeTaken, _ := app.callGetEmail(session)
		session.APIsTimeTaken.SetEmailAPITime(timeTaken)
	}()

	wg.Wait()
}

func (app *App) callGetReport(session *Session) (string, int64, error) {
	if session == nil {
		return "", 0, fmt.Errorf("sesssion should be non-nil value")
	}
	app.InfoLogger.Println("Sending Request to get report for session ID:", session.ID)
	reportStart := time.Now()
	report, err := app.QuizAPI.GetReport(session.ID)
	reportEnd := time.Now()
	app.InfoLogger.Printf("Report received for session ID: %s, report: %+v\n", session.ID, report)
	if err != nil {
		session.SetError(err)
		session.SetStatus(STATUS_FAILED)
		session.SetEndTime(time.Now())
		app.ErrorLogger.Printf("Error getting report for session ID: %s, error: %v\n", session.ID, err)
		app.Errors <- &SessionError{
			Session: session,
		}
		return "", reportEnd.UnixMilli() - reportStart.UnixMilli(), err
	}
	return report, reportEnd.UnixMilli() - reportStart.UnixMilli(), nil
}

func (app *App) callGetEmail(session *Session) (int64, error) {
	if session == nil {
		return 0, fmt.Errorf("sesssion should be non-nil value")
	}
	app.InfoLogger.Println("Sending Request to get email report for session ID:", session.ID)
	emailStart := time.Now()
	_, err := app.QuizAPI.GetEmailReport(session.ID)
	emailEnd := time.Now()
	app.InfoLogger.Printf("Email Request Successful for session ID: %s\n", session.ID)
	if err != nil {
		session.SetError(err)
		session.SetStatus(STATUS_FAILED)
		session.SetEndTime(emailEnd)
		app.Errors <- &SessionError{
			Session: session,
		}
		return emailEnd.UnixMilli() - emailStart.UnixMilli(), err
	}
	return emailEnd.UnixMilli() - emailStart.UnixMilli(), nil
}

// NOTE: FakeUserAction can be used instead of SimulateUser to run the application
//
// func (app *App) FakeUserAction(email, topic string) {
// 	defer app.Wait.Done()
// 	log.Printf("Simulating user action for email: %s, topic: %s\n", email, topic)
// 	startTime := time.Now()
// 	time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
// 	endTime := time.Now()
//
// 	app.Results <- &Session{
// 		ID:        "somereandomstring",
// 		UserID:    "user_" + email,
// 		Email:     email,
// 		Topic:     topic,
// 		Status:    STATUS_COMPLETED,
// 		StartTime: startTime.UnixMilli(),
// 		EndTime:   endTime.UnixMilli(),
// 	}
// }
