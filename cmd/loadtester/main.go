package main

import (
	"runtime"
	"sync"
	"time"
)

func main() {
	// set the number of go routines
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	// Number of users to simulate
	numUsers := 1000

	// Create a wait group to wait for all user simulations to complete
	wg := sync.WaitGroup{}

	// Start the user simulation
	for i := range numUsers {
		wg.Add(1)
		go func(userID int) {
			defer func() {
				if r := recover(); r != nil {
					println("Recovered from panic in user", userID, ":", r)
				}
				wg.Done()
			}()
			// Simulate user activity
			time.Sleep(time.Second * 2)
			// TODO: Add simulation logic here
		}(i)
	}

	// Wait for the simulation to complete
	wg.Wait()
}
