package main

import (
	"runtime"
	"time"

	"github.com/go-squad-5/quiz-load-test/internal/app"
)

func main() {
	// set the number of go routines
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	// Number of users to simulate
	numUsers := 1000

  // application
  application := app.NewApp()

	// Start the user simulation
	for i := range numUsers {
		application.Wait.Add(1)
		go func(userID int) {
			defer func() {
				if r := recover(); r != nil {
					println("Recovered from panic in user", userID, ":", r)
				}
				application.Wait.Done()
			}()
			// TODO: Add simulation logic here
			time.Sleep(time.Second * 2)
		}(i)
	}

	// Wait for the simulation to complete
	application.Wait.Wait()
}
