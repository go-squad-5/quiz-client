package app

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_app_config_LoadConfig_WhenNoEnvs(t *testing.T) {
	// should set default values in the config struct
	config := LoadConfig()
	require.NotNil(t, config, "Expected returned value to be non-nil, but got nil value")

	defaultBaseUrl := "http://localhost:8080"
	defaultReportServerUrl := "http://localhost:8070"
	defaultNumUsers := 10
	assert.Equalf(t, defaultBaseUrl, config.BaseURL, "Expected base url to be default value %s, but got %s", defaultBaseUrl, config.BaseURL)
	assert.Equalf(t, defaultReportServerUrl, config.ReportServerBaseURL, "Expected report server base url to be default value %s, but got %s", defaultReportServerUrl, config.ReportServerBaseURL)
	assert.Equalf(t, defaultNumUsers, config.NumUsers, "Expected default number of users to be %d, but got %d", defaultNumUsers)
}

func Test_app_config_LoadConfig_WhenSetEnvs(t *testing.T) {
	baseUrl := "http://localhost:3000"
	reportServerUrl := "http://localhost:3070"
	numUsers := 100
	os.Setenv("BASE_URL", baseUrl)
	os.Setenv("REPORT_SERVER_BASEURL", reportServerUrl)
	os.Setenv("NUM_USERS", strconv.Itoa(numUsers))

	// should set os env values in the config struct
	config := LoadConfig()
	require.NotNil(t, config, "Expected returned value to be non-nil, but got nil value")

	assert.Equalf(t, baseUrl, config.BaseURL, "Expected base url to be  value %s, but got %s", baseUrl, config.BaseURL)
	assert.Equalf(t, reportServerUrl, config.ReportServerBaseURL, "Expected report server base url to be  value %s, but got %s", reportServerUrl, config.ReportServerBaseURL)
	assert.Equalf(t, numUsers, config.NumUsers, "Expected  number of users to be %d, but got %d", numUsers)
}

func Test_app_config_LoadConfig_WhenURLsWithTrailingSlash(t *testing.T) {
	baseUrl := "http://localhost:3000"
	reportServerUrl := "http://localhost:3070"
	numUsers := 100
	os.Setenv("BASE_URL", baseUrl+"/")
	os.Setenv("REPORT_SERVER_BASEURL", reportServerUrl+"/")
	os.Setenv("NUM_USERS", strconv.Itoa(numUsers))

	// should set os env values in the config struct
	config := LoadConfig()
	require.NotNil(t, config, "Expected returned value to be non-nil, but got nil value")

	assert.Equalf(t, baseUrl, config.BaseURL, "Expected base url to be  value %s, but got %s", baseUrl, config.BaseURL)
	assert.Equalf(t, reportServerUrl, config.ReportServerBaseURL, "Expected report server base url to be  value %s, but got %s", reportServerUrl, config.ReportServerBaseURL)
	assert.Equalf(t, numUsers, config.NumUsers, "Expected  number of users to be %d, but got %d", numUsers)
}

func Test_app_config_LoadConfig_WhenInvalidNumUsers(t *testing.T) {
	baseUrl := "http://localhost:3000"
	reportServerUrl := "http://localhost:3070"
	numUsers := "invalid"
	os.Setenv("BASE_URL", baseUrl)
	os.Setenv("REPORT_SERVER_BASEURL", reportServerUrl)
	os.Setenv("NUM_USERS", numUsers)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected LoadConfig to panic with invalid number of users, but it did not")
		}
	}()

	LoadConfig()
}
