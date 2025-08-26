package quizapi

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_quizapi_report_BuildGetReportAPIURL(t *testing.T) {
	sessionID := "12345"
	expectedURL := "http://example.com/sessions/12345/report"
	actualURL := buildGetReportAPIURL("http://example.com/sessions/%s/report", sessionID)
	assert.Equal(t, expectedURL, actualURL, "URLs should match")
}

func Test_quizapi_report_ValidateGetReportResponseStatus(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError bool
	}{
		{"success status", 200, false},
		{"bad request", 400, true},
		{"unauthorized", 401, true},
		{"not found", 404, true},
		{"internal server error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			err := validateGetReportResponseStatus(resp)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error for status code %d", tt.statusCode)
			} else {
				assert.NoError(t, err, "Did not expect an error for status code %d", tt.statusCode)
			}
		})
	}
}

func Test_quizapi_report_ParseGetReportErrorResponse(t *testing.T) {
	type GetReportErrorResponse struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	}
	tests := []struct {
		name          string
		body          string
		message       string
		statusCode    int
		expectedError bool
	}{
		{"valid error json", `{"message": "error processing request", "statusCode": 400}`, "error processing request", 400, false},
		{"invalid error json", `{"message": "error processing request}`, "error processing request", 0, true},
		{"missing status code", `{"message": "error processing request"}`, "error processing request", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tt.body))
			errResp, err := parseGetReportErrorResponse(body)
			if tt.expectedError {
				require.Error(t, err, "Expected an error while parsing response, got nil error")
				return
			}
			require.NoError(t, err, "Expected no error to be returned, got an error %w", err)
			assert.Contains(t, errResp, tt.message)
			assert.Contains(t, errResp, strconv.Itoa(tt.statusCode))
		})
	}
}

func Test_quizapi_report_OpenSessionReportFile_WhenDirExist(t *testing.T) {
	reportsDirPath = "../../tmp/reports" // change path for the test environment
	sessionId := "test_session"
	expectedFilePath := fmt.Sprintf("%s/%s_report.pdf", reportsDirPath, sessionId)

	defer func() {
		if err := os.Remove(expectedFilePath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Error cleaning up the test file: %s", expectedFilePath)
		}
	}()

	file, filePath, err := openSessionReportFile(sessionId)
	require.NoErrorf(t, err, "Expected no error while opening a file, but got error %w", err)
	defer file.Close()

	assert.NotNil(t, file)
	assert.Equalf(t, expectedFilePath, filePath, "Expected filePath to be correct path in %s/%s_reports.pdf", reportsDirPath, sessionId)
}

func Test_quizapi_report_OpenSessionReportFile_WhenDirNotExist(t *testing.T) {
	reportsDirPath = "../../tmp/reports" // change path for the test environment
	sessionId := "test_session"
	expectedFilePath := fmt.Sprintf("%s/%s_report.pdf", reportsDirPath, sessionId)

	// remove the reports directory
	if err := os.RemoveAll(reportsDirPath); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove the %s directory", reportsDirPath)
	}

	defer func() {
		if err := os.Remove(expectedFilePath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Error cleaning up the test file: %s", expectedFilePath)
		}
	}()

	file, filePath, err := openSessionReportFile(sessionId)
	require.NoErrorf(t, err, "Expected no error while opening a file, but got error %w", err)
	defer file.Close()

	assert.NotNil(t, file)
	assert.Equalf(t, expectedFilePath, filePath, "Expected filePath to be correct path in %s/%s_reports.pdf", reportsDirPath, sessionId)
}

func readFile(r *os.File) (string, error) {
	res := ""
	buf := make([]byte, 1)
	for {
		_, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		res += string(buf)
	}
	return res, nil
}

func Test_quizapi_report_SaveResponseToFile_ValidFile(t *testing.T) {
	data := "Hello World"
	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(data)),
	}

	// mock file to get the data written
	t.Log("Creating a os pipe")
	r, w, err := os.Pipe()
	require.NoError(t, err, "Failed to create a read-write pipe for testing")
	defer r.Close()
	defer w.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = saveResponseToFile(response, w)
		w.Close()
		require.NoErrorf(t, err, "Expected no error while saving response to file, but got error: %w", err)
	}()

	// read data from file
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := readFile(r)
		require.NoErrorf(t, err, "Failed to read data from file.")
		assert.Equalf(t, data, res, "Expected written data to the file to be %s, but got %s", data, res)
	}()

	wg.Wait()
}

func Test_quizapi_report_SaveResponseToFile_WhenClosedRespBody(t *testing.T) {
	data := "Hello World"
	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(data)),
	}
	response.Body.Close() // close the body early

	// mock file to get the data written
	t.Log("Creating a os pipe")
	r, w, err := os.Pipe()
	require.NoError(t, err, "Failed to create a read-write pipe for testing")
	defer r.Close()
	w.Close() // close early

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = saveResponseToFile(response, w)
		w.Close()
		require.Errorf(t, err, "Expected an error while saving to a closed file, but got no error")
		assert.Containsf(t, err.Error(), "failed to write file", "Expected error to contain 'failed to write file', got: %s", err.Error())
		assert.Containsf(t, err.Error(), "closed", "Expected error to contain 'closed', got: %s", err.Error())
	}()

	// read data from file
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := readFile(r)
		require.NoErrorf(t, err, "Failed to read data from file.")
		assert.NotEqualf(t, data, res, "Expected data to not being written, passed data: %s, written data: %s", data, res)
	}()

	wg.Wait()
}

func Test_quizapi_report_SaveResponseToFile_WhenClosedFile(t *testing.T) {
	data := "Hello World"
	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(data)),
	}

	// mock file to get the data written
	t.Log("Creating a os pipe")
	r, w, err := os.Pipe()
	require.NoError(t, err, "Failed to create a read-write pipe for testing")
	defer r.Close()
	w.Close() // close early

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = saveResponseToFile(response, w)
		w.Close()
		require.Errorf(t, err, "Expected an error while saving to a closed file, but got no error")
		assert.Containsf(t, err.Error(), "failed to write file", "Expected error to contain 'failed to write file', got: %s", err.Error())
		assert.Containsf(t, err.Error(), "closed", "Expected error to contain 'closed pipe', got: %s", err.Error())
	}()

	// read data from file
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := readFile(r)
		require.NoErrorf(t, err, "Failed to read data from file.")
		assert.NotEqualf(t, data, res, "Expected data to not being written, passed data: %s, written data: %s", data, res)
	}()

	wg.Wait()
}

func Test_quizapi_report_ParseJsonGetReportResponse(t *testing.T) {
	tests := []struct {
		name          string
		data          GetReportResponse
		json          string
		expectedError bool
	}{
		{
			"valid json",
			GetReportResponse{true, "success", GetReportResponseData{"1", "abc", "url", "12", "content"}},
			`{"success":true,"message":"success","data":{"documentId":"1","fileName":"abc","downloadUrl":"url","expiresAt":"12","contentBase64":"content"}}`,
			false,
		},
		{
			"invalid json",
			GetReportResponse{true, "success", GetReportResponseData{"1", "abc", "url", "12", "content"}},
			`{"success":true,"message":"success","data":{"documentId":"1","fileName":"abc","downloadUrl":"url","expiresAt":"12","contentBase64":"content"}`,
			true,
		},
		{
			"missing data",
			GetReportResponse{true, "success", GetReportResponseData{}},
			`{"success":true,"message":"success","data":{}}`,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tt.json))
			resp, err := parseJsonGetReportResponse(body)
			if tt.expectedError {
				require.Error(t, err, "Expected an error while parsing response, got nil error")
				return
			}
			require.NoError(t, err, "Expected no error to be returned, got an error %w", err)
			assert.Equal(t, tt.data.Success, resp.Success)
			assert.Equal(t, tt.data.Message, resp.Message)
			assert.Equal(t, tt.data.Data.DocumentID, resp.Data.DocumentID)
			assert.Equal(t, tt.data.Data.FileName, resp.Data.FileName)
			assert.Equal(t, tt.data.Data.DownloadURL, resp.Data.DownloadURL)
			assert.Equal(t, tt.data.Data.ExpiresAt, resp.Data.ExpiresAt)
			assert.Equal(t, tt.data.Data.ContentBase64, resp.Data.ContentBase64)
		})

	}
}

func Test_quizapi_report_DecodeBase64Content_ValidContent(t *testing.T) {
	data := "Hello World"
	base64Content := base64.StdEncoding.EncodeToString([]byte(data))

	// mock file to get the data written
	r, w, err := os.Pipe()
	require.NoError(t, err, "Failed to create a read-write pipe for testing")
	defer r.Close()
	defer w.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = decodeAndSaveBase64Response(base64Content, w)
		w.Close()
		require.NoErrorf(t, err, "Expected no error while saving response to file, but got error: %w", err)
	}()

	// read data from file
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := readFile(r)
		require.NoErrorf(t, err, "Failed to read data from file.")
		assert.Equalf(t, data, res, "Expected written data to the file to be %s, but got %s", data, res)
	}()

	wg.Wait()
}

func Test_quizapi_report_DecodeBase64Content_InvalidContent(t *testing.T) {
	data := "Hello World"
	// some invalid base64 content
	base64Content := "invalid base64 content"

	// mock file to get the data written
	t.Log("Creating a os pipe")
	r, w, err := os.Pipe()
	require.NoError(t, err, "Failed to create a read-write pipe for testing")
	defer r.Close()
	defer w.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = decodeAndSaveBase64Response(base64Content, w)
		w.Close()
		require.Errorf(t, err, "Expected an error while saving invalid response to file, but got no error")
	}()

	// read data from file
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := readFile(r)
		require.NoErrorf(t, err, "Failed to read data from file.")
		assert.NotEqualf(t, data, res, "Expected data to not being written, passed data: %s, written data: %s", data, res)
	}()

	wg.Wait()
}

func Test_quizapi_report_DecodeBase64Content_WhenClosedFile(t *testing.T) {
	data := "Hello World"
	base64Content := base64.StdEncoding.EncodeToString([]byte(data))

	// mock file to get the data written
	t.Log("Creating a os pipe")
	r, w, err := os.Pipe()
	require.NoError(t, err, "Failed to create a read-write pipe for testing")
	defer r.Close()
	w.Close() // close early

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = decodeAndSaveBase64Response(base64Content, w)
		w.Close()
		require.Errorf(t, err, "Expected an error while saving to a closed file, but got no error")
		assert.Containsf(t, err.Error(), "failed to write file", "Expected error to contain 'failed to write file', got: %s", err.Error())
		assert.Containsf(t, err.Error(), "closed", "Expected error to contain 'closed', got: %s", err.Error())
	}()

	// read data from file
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := readFile(r)
		require.NoErrorf(t, err, "Failed to read data from file.")
		assert.NotEqualf(t, data, res, "Expected data to not being written, passed data: %s, written data: %s", data, res)
	}()

	wg.Wait()
}

func Test_quizapi_report_GetReport_WhenSuccess(t *testing.T) {
	reportsDirPath = "../../tmp/reports"
	baseUrl := "http://localhost:3000"
	reportServerBaseUrl := "http://localhost:3001"
	ssid := "12345"

	q := NewTestQuizAPI(baseUrl, reportServerBaseUrl, func(req *http.Request) *http.Response {
		assert.Equalf(t, http.MethodGet, req.Method, "Expected GET request, but got %s", req.Method)
		expectedUrl := fmt.Sprintf("%s/sessions/%s/report", reportServerBaseUrl, ssid)
		actualUrl := req.URL.String()
		assert.Equalf(t, expectedUrl, actualUrl, "Expected request URL to match the report endpoint %s, but got %s", expectedUrl, actualUrl)

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("Hello World")),
		}
	})

	filePath, err := q.GetReport(ssid)
	require.NoError(t, err, "Expected no error while getting report, but got error %v", err)

	assert.NotEmpty(t, filePath, "Expected file path to be returned, but got empty string")
	assert.Equalf(t, filePath, fmt.Sprintf("%s/%s_report.pdf", reportsDirPath, ssid), "Expected report file path to match %s/%s_report.pdf, but got %s", reportsDirPath, ssid, filePath)
}

func Test_quizapi_report_GetReport_WhenError(t *testing.T) {
	reportsDirPath = "../../tmp/reports"
	baseUrl := "http://localhost:3000"
	reportServerBaseUrl := "http://localhost:3001"
	ssid := "12345"

	q := NewTestQuizAPI(baseUrl, reportServerBaseUrl, func(req *http.Request) *http.Response {
		assert.Equalf(t, http.MethodGet, req.Method, "Expected GET request, but got %s", req.Method)
		expectedUrl := fmt.Sprintf("%s/sessions/%s/report", reportServerBaseUrl, ssid)
		actualUrl := req.URL.String()
		assert.Equalf(t, expectedUrl, actualUrl, "Expected request URL to match the report endpoint %s, but got %s", expectedUrl, actualUrl)

		return &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(strings.NewReader(`{"statusCode": 400, "message": "error processing request"}`)),
		}
	})

	filePath, err := q.GetReport(ssid)
	require.Error(t, err, "Expected an error while getting report, but got no error")

	assert.Empty(t, filePath, "Expected file path to be empty, but got %s", filePath)
	assert.Contains(t, err.Error(), "400", "Expected error message to contain 'status code 400', but got %s", err.Error())
	assert.Contains(t, err.Error(), "error processing request", "Expected error message to contain 'error processing request', but got %s", err.Error())
}

func Test_quizapi_report_GetReport_WhenErrorWithInvalidResponse(t *testing.T) {
	reportsDirPath = "../../tmp/reports"
	baseUrl := "http://localhost:3000"
	reportServerBaseUrl := "http://localhost:3001"
	ssid := "12345"

	q := NewTestQuizAPI(baseUrl, reportServerBaseUrl, func(req *http.Request) *http.Response {
		assert.Equalf(t, http.MethodGet, req.Method, "Expected GET request, but got %s", req.Method)
		expectedUrl := fmt.Sprintf("%s/sessions/%s/report", reportServerBaseUrl, ssid)
		actualUrl := req.URL.String()
		assert.Equalf(t, expectedUrl, actualUrl, "Expected request URL to match the report endpoint %s, but got %s", expectedUrl, actualUrl)

		return &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(strings.NewReader(`{"statusCode": 400, "message": "error processing request`)),
		}
	})

	filePath, err := q.GetReport(ssid)
	require.Error(t, err, "Expected an error while getting report, but got no error")

	assert.Empty(t, filePath, "Expected file path to be empty, but got %s", filePath)
	assert.Contains(t, err.Error(), "failed to parse error response", "Expected error message to contain 'failed to parse error response', but got %s", err.Error())
}

func Test_quizapi_report_GetReport_WhenNetworkError(t *testing.T) {
	reportsDirPath = "../../tmp/reports"
	baseUrl := "http://localhost:3000"
	reportServerBaseUrl := "http://localhost:3001"
	ssid := "12345"
	q := NewTestQuizAPI(baseUrl, reportServerBaseUrl, func(req *http.Request) *http.Response {
		return nil
	})

	_, err := q.GetReport(ssid)
	assert.Error(t, err, "Expected an network error while getting report, but got no error")
}
