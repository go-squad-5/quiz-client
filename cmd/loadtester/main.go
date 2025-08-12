package main

import (
	"runtime"
	"time"

	application "github.com/go-squad-5/quiz-load-test/internal/app"
)

func main() {
	// set the number of os threads to use for the simulation
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	startTime := time.Now()

	app := application.NewApp()

	app.ErrorListener.Add(1)
	app.InfoLogger.Println("GO ROUTINE Started for listening to errors")
	go app.ListenForErrors()

	app.ResultListener.Add(1)
	app.InfoLogger.Println("GO ROUTINE Started for listening to results")
	go app.ListenForResults()

	app.InfoLogger.Println("Starting simulation with", app.Config.NumUsers, "users")
	app.StartSimulation()

	app.Wait.Wait()
	elapsed := time.Since(startTime)

	app.Stop()
	elapsed2 := time.Since(startTime)

	app.ResultLogger.Println(
		"Total time taken to complete all sessions concurrently: ",
		elapsed.Seconds(),
		" seconds",
	)
	app.ResultLogger.Println(
		"Total time taken by test: ",
		elapsed2.Seconds(),
		" seconds",
	)
}
