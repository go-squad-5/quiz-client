# Quiz Client

The Quiz client contains the quiz CLI and the quiz load tester.

# Quiz CLI (`cmd/cli`)
The quiz CLI is a command line interface for the quiz server. It allows you to start a quiz, answer questions, and view results.

# Quiz Load Tester (`cmd/loadtester`)
The quiz load tester is a tool for testing the performance of the quiz server. It can simulate multiple users answering questions and submitting results.

## Setup Environment Variables
To run the quiz CLI and load tester, you need to set up the environment variables, example declared in `.env.example`:
```bash
cp .env.example .env
```
Then, edit the `.env` file to set your environment variables.

## Run Load Tester
To run the load tester, you need to have the quiz server running. You can start the load tester with the following command:

```bash
go run ./cmd/loadtester
```

> Check the logs from the `./tmp/logs.txt` file
> Check Quiz Reports for each session in the `./tmp/reports` directory
