package app

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func (app *App) ListenForResults() {
	defer app.ResultListener.Done()
	defer app.InfoLogger.Println("GO ROUTINE FINISHED for listening to results")

	file := openResultsFile()
	defer file.Close()

	timetaken := []int64{}
	// listen for results from the simulation and log them into the file
	for result := range app.Results {
		// get the result log string
		logString := getResultLog(result)

		// aggregate the results time taken
		timetaken = append(timetaken, result.EndTime-result.StartTime)

		// write the log string to the file
		_, err := file.WriteString(logString)
		if err != nil {
			panic("Failed to write to results file: " + err.Error())
		}
	}

	summary := getSummaryLog(timetaken, app.Config.NumUsers)
	fmt.Print(summary)

	// write the summary to the file
	_, err := file.WriteString(summary)
	if err != nil {
		panic("Failed to write summary to results file: " + err.Error())
	}
}

func openResultsFile() *os.File {
	// create a tmp directory if it doesn't exist
	if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
		err := os.Mkdir("./tmp", 0755)
		if err != nil {
			panic("Failed to create tmp directory: " + err.Error())
		}
	}
	// open the results file
	file, err := os.Create("./tmp/logs.txt")
	if err != nil {
		panic("Failed to create results file: " + err.Error())
	}
	return file
}

func getResultLog(result *Session) string {
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
	logString = logString + "Start Time: " + time.UnixMilli(result.StartTime).Format(time.RFC3339) + "\n"
	logString = logString + "End Time: " + time.UnixMilli(result.EndTime).Format(time.RFC3339) + "\n"
	logString = logString + "Time Taken: " + strconv.FormatInt(result.EndTime-result.StartTime, 10) + " ms\n"
	logString = logString + "Questions-Answers: " + fmt.Sprintf("%v", result.Answers) + "\n"
	logString = logString + "Report: " + result.Report + "\n"
	if result.Error != nil {
		logString = logString + "Error: " + result.Error.Error() + "\n"
	}
	if result.APIsTimeTaken != nil {
		logString = logString + "APIs Time Taken:\n"
		logString = logString + "Session Creation: " + strconv.FormatInt(result.APIsTimeTaken.SessionCreation, 10) + " ms\n"
		logString = logString + "Start Quiz: " + strconv.FormatInt(result.APIsTimeTaken.StartQuiz, 10) + " ms\n"
		logString = logString + "Submit Quiz: " + strconv.FormatInt(result.APIsTimeTaken.SubmitQuiz, 10) + " ms\n"
		logString = logString + "Report API: " + strconv.FormatInt(result.APIsTimeTaken.ReportAPI, 10) + " ms\n"
		logString = logString + "Email API: " + strconv.FormatInt(result.APIsTimeTaken.EmailAPI, 10) + " ms\n"
	} else {
		logString = logString + "APIs Time Taken: Not available\n"
	}
	logString = logString + "-----------------------------------------------\n"
	return logString
}

func getSummaryLog(timetaken []int64, numOfUsers int) string {
	var totalTime int
	for _, time := range timetaken {
		totalTime += int(time)
	}

	averageTime := float64(totalTime) / float64(len(timetaken))

	summary := "-------------------RESULTS--------------------\n"
	summary += "Total Sessions: " + strconv.Itoa(numOfUsers) + "\n"
	summary += "Average Time Taken per session: " + strconv.FormatFloat(averageTime, 'f', 2, 64) + " milliseconds\n"
	summary += "Check ./tmp/logs.txt for all logs\n"
	summary += "-----------------------------------------------\n"

	return summary
}
