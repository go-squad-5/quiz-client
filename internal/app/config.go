package app

import (
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	BaseURL             string
	ReportServerBaseURL string
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
		baseUrl = "http://localhost:3000"
	}
	reportServerBaseUrl := os.Getenv("REPORT_SERVER_BASEURL")
	if reportServerBaseUrl == "" {
		reportServerBaseUrl = "http://localhost:3002"
	}

	// trim trailing slashes
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	reportServerBaseUrl = strings.TrimSuffix(reportServerBaseUrl, "/")

	return &Config{
		BaseURL:             baseUrl,
		ReportServerBaseURL: reportServerBaseUrl,
	}
}
