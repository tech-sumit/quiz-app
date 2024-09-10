# Quiz App API

This is a simple Quiz App API built with Go, following the MVC architecture.

## Requirements

- Go 1.16 or higher
- Gorilla Mux package

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/tech-sumit/quiz-app.git
   ```

2. Navigate to the project directory:
   ```
   cd quiz-app
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

## Running the Application

To run the application, use the following command:

```
go run cmd/server/main.go
```

The server will start running on `http://localhost:8080`.

## Running Tests

To run the tests, use the following command:

```
go test ./
