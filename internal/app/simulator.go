package app

import (
	"fmt"
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
		go app.SimulateUser(email, topic)
	}
}

func (app *App) SimulateUser(email, topic string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in user", email, ":", r)
		}
		app.Wait.Done()
	}()

	fmt.Printf("Simulating user action for email: %s, topic: %s\n", email, topic)

	startTime := time.Now()

	// -------------------------------------
	// // Fake Simulation of user action
	// app.FakeUserAction(email, topic)
	// endTime := time.Now()
	//
	// app.Results <- &Session{
	// 	ID:        "somereandomstring",
	// 	UserID:    "user_" + email,
	// 	Email:     email,
	// 	Topic:     topic,
	// 	Status:    STATUS_COMPLETED,
	// 	StartTime: startTime.Unix(),
	// 	EndTime:   endTime.Unix(),
	// }
	// return
	// -------------------------------------

	// Create a session
	ssid, err := app.QuizAPI.CreateSession(email, topic)
	if err != nil {
		app.Errors <- &StartSessionError{
			Email: email,
			Topic: topic,
			err:   err,
		}
		return
	}

	session := &Session{
		ID:        ssid,
		UserID:    "user_" + email,
		Email:     email,
		Topic:     topic,
		Status:    STATUS_STARTED,
		StartTime: startTime.Unix(),
		EndTime:   0,
		Answers:   []quizapi.Answer{},
		Score:     0,
		CreatedAt: startTime.Unix(),
	}

	// Start the quiz , get questions
	questions, err := app.QuizAPI.StartQuiz(ssid)
	if err != nil {
		session.Error = err
		session.Status = STATUS_FAILED
		session.EndTime = time.Now().Unix()
		app.Errors <- &StartSessionError{
			Email: email,
			Topic: topic,
			err:   err,
		}
	}

	session.Question = questions

	// mark random answers to questions
	for _, question := range questions {
		numOptions := len(question.Options)
		if numOptions == 0 {
			session.Error = fmt.Errorf("no options available for question ID: %s", question.ID)
			session.Status = STATUS_FAILED
			session.EndTime = time.Now().Unix()
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
	score, err := app.QuizAPI.SubmitQuiz(ssid, session.Answers)
	if err != nil {
		session.Status = STATUS_FAILED
		session.Error = err
		session.EndTime = time.Now().Unix()
		app.Errors <- &SessionError{
			Session: session,
		}
		return
	}

	// update session with score and end time
	session.Score = score

	// Get the report for the session
	report, err := app.QuizAPI.GetReport(ssid, session.UserID) // TODO: check if userID is needed
	if err != nil {
		session.Error = err
		session.Status = STATUS_FAILED
		session.EndTime = time.Now().Unix()
		app.Errors <- &SessionError{
			Session: session,
		}
		return
	}
	session.Report = report

	endTime := time.Now()
	session.EndTime = endTime.Unix()
	session.Status = STATUS_COMPLETED

	app.Results <- session
}

func (app *App) FakeUserAction(email, topic string) {
	time.Sleep(1 * time.Second)
}
