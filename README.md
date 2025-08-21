# Quiz Client

The Quiz client contains the quiz client to load test services.

# Quiz Load Tester (`cmd/loadtester`)
The quiz load tester is a tool for testing the performance of the quiz server. It can simulate multiple users answering questions and submitting results.

## Setup
- To run the quiz CLI and load tester, you need to set up the environment variables, example declared in `.env.example`:
```bash
cp .env.example .env
```
Then, edit the `.env` file to set your environment variables.

- To install the required dependencies, you can use the following command:
```bash
go mod tidy
```

- To build the quiz client, you can use the following command:
```bash
go build -o ./bin/quiz ./cmd/loadtester
```
Then you can run the quiz client: `./bin/quiz`

## Run Load Tester
To run the load tester, you need to have the quiz server running. You can start the load tester with the following command:

```bash
go run ./cmd/loadtester
```

> In order to set number of users to simulate, set `NUM_USERS` environment variable, defaults to 10, defaults to 10.
> Check the logs from the `./tmp/logs.txt` file
> Check Quiz Reports for each session in the `./tmp/reports` directory

## Run Tests

- To run the tests for the quiz client, you can use the following command:

```bash
go test -cover ./...
```

- For better test coverage, you can run the tests with the following command:

```bash
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

- Using `tparse` you can parse the test results and generate a report:

```bash
go test -cover -coverprofile=coverage.out -json ./... | tparse -all
```
> to install tparse: `go install github.com/mfridman/tparse@latest`

### Running Services:-
- Session Manager
- Quiz Master
- PDF Generator (Report and Email Service)
> Set the respective urls with ports in the service's envs in order to integrate all services


## Test Coverage Report

| Status | Elapsed | Package                                                | Cover | Pass | Fail | Skip |
|--------|---------|--------------------------------------------------------|-------|------|------|------|
| PASS | 0.01s  | github.com/go-squad-5/quiz-load-test/cmd/loadtester     | 0.0%  | 0    | 0    | 0    |
| PASS | 1.93s  | github.com/go-squad-5/quiz-load-test/internal/app       | 96.9% | 61   | 0    | 0    |
| PASS | 0.02s  | github.com/go-squad-5/quiz-load-test/internal/quizapi   | 94.0% | 113  | 0    | 0    |
| PASS | 0.01s  | github.com/go-squad-5/quiz-load-test/internal/quizapi/mock | 0.0%  | 0    | 0    | 0    |


## Detailed Coverage Report 

Checkout detailed coverage report at [https://go-squad-5.github.io/quiz-client/coverage](https://go-squad-5.github.io/quiz-client/coverage#file1)
