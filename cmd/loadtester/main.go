package main

import (
	"fmt"
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
	go app.ListenForErrors()

	app.ResultListener.Add(1)
	go app.ListenForResults()

	app.StartSimulation()

	app.Wait.Wait()
	elapsed := time.Since(startTime)

	app.Stop()
	elapsed2 := time.Since(startTime)

	fmt.Println(
		"Total time taken to complete all sessions concurrently: ",
		elapsed.Seconds(),
		" seconds",
	)
	fmt.Println(
		"Total time taken by test: ",
		elapsed2.Seconds(),
		" seconds",
	)
}
