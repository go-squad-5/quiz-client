package app

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func (app *App) ListenForResults() {
	defer app.Listeners.Done()

	file := openResultsFile()
	defer file.Close()

	timetaken := []int64{}
	// listen for results from the simulation and log them into the file
	for result := range app.Results {
		logString := ""
		ssid := result.ID
		if ssid == "" {
			logString = logString + "----------------Error while starting------------\n"
		} else {
			logString = logString + "----------------Session ID: " + ssid + "----------------\n"
		}
		logString = logString + "Email: " + result.Email + "\n"
		logString = logString + "User ID: " + result.UserID + "\n"
		logString = logString + "Score: " + strconv.Itoa(result.Score) + "\n"
		logString = logString + "Status: " + string(result.Status) + "\n"
		logString = logString + "Start Time: " + time.Unix(result.StartTime, 0).Format(time.RFC3339) + "\n"
		logString = logString + "End Time: " + time.Unix(result.EndTime, 0).Format(time.RFC3339) + "\n"
		logString = logString + "Time Taken: " + strconv.FormatFloat(float64(result.EndTime-result.StartTime)/1000, 'f', 2, 64) + " seconds\n"
		logString = logString + "Questions-Answers: " + fmt.Sprintf("%v", result.Answers) + "\n"
		if result.Error != nil {
			logString = logString + "Error: " + result.Error.Error() + "\n"
		}
		logString = logString + "-----------------------------------------------\n"

		// aggregate the results
		timetaken = append(timetaken, result.EndTime-result.StartTime)

		// write the log string to the file
		_, err := file.WriteString(logString)
		if err != nil {
			panic("Failed to write to results file: " + err.Error())
		}
	}

	// calculate the average time taken
	var total int
	for _, time := range timetaken {
		total += int(time)
	}

	averageTime := float64(total) / float64(len(timetaken))
	summary := "-----------------------------------------------\n"
	summary += "Total Sessions: " + strconv.Itoa(app.Config.NumUsers) + "\n"
	summary += "Average Time Taken per session: " + strconv.FormatFloat(averageTime, 'f', 2, 64) + " seconds\n"
	summary += "Check ./tmp/results.txt for all logs\n"
	summary += "-----------------------------------------------\n"
	fmt.Print(summary)

	// write the summary to the file
	_, err := file.WriteString(summary)
	if err != nil {
		panic("Failed to write summary to results file: " + err.Error())
	}
}

func openResultsFile() *os.File {
	// open the results file
	// create a tmp directory if it doesn't exist
	if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
		err := os.Mkdir("./tmp", 0755)
		if err != nil {
			panic("Failed to create tmp directory: " + err.Error())
		}
	}
	file, err := os.Create("./tmp/results.txt")
	if err != nil {
		panic("Failed to create results file: " + err.Error())
	}
	return file
}
