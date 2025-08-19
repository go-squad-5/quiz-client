package app

import (
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	BaseURL             string
	ReportServerBaseURL string
	NumUsers            int
}

type Endpoints struct {
	CreateSession string
	StartQuiz     string
	SubmitQuiz    string
	GetReport     string
}

func LoadConfig() *Config {
	// load configurations from environment variables
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:8080"
	}
	reportServerBaseUrl := os.Getenv("REPORT_SERVER_BASEURL")
	if reportServerBaseUrl == "" {
		reportServerBaseUrl = "http://localhost:8070"
	}
	numUsers := os.Getenv("NUM_USERS")
	if numUsers == "" {
		numUsers = "10"
	}

	// trim trailing slashes
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	reportServerBaseUrl = strings.TrimSuffix(reportServerBaseUrl, "/")

	// convert numUsers to int
	numUsersInt, err := strconv.Atoi(numUsers)
	if err != nil {
		panic("Invalid NUM_USERS value, must be an integer")
	}

	return &Config{
		BaseURL:             baseUrl,
		ReportServerBaseURL: reportServerBaseUrl,
		NumUsers:            numUsersInt,
	}
}
