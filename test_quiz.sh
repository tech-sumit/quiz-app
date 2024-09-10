#!/bin/bash

# Function to run curl command and display output
run_curl() {
    echo "----------------------------------------------"
    echo "Running: $1"
    echo "----------------------------------------------"
    eval $1
    echo -e "\n"
}

# Set the base URL
BASE_URL="http://localhost:8080"

# 1. Create a Quiz
CREATE_QUIZ_CMD="curl -s -X POST $BASE_URL/quiz \
  -H \"Content-Type: application/json\" \
  -d '{
    \"id\": \"1\",
    \"title\": \"Test Quiz\",
    \"is_negative_marking\": true,
    \"penalty\": 0.5,
    \"questions\": [
      {
        \"id\": \"q1\",
        \"text\": \"What is 1+1?\",
        \"options\": [\"1\", \"2\", \"3\", \"4\"],
        \"correct_option\": 1,
        \"marks\": 2
      },
      {
        \"id\": \"q2\",
        \"text\": \"What is the capital of France?\",
        \"options\": [\"London\", \"Berlin\", \"Paris\", \"Madrid\"],
        \"correct_option\": 2,
        \"marks\": 3
      }
    ]
  }'"

run_curl "$CREATE_QUIZ_CMD"

# 2. Retrieve the Created Quiz
GET_QUIZ_CMD="curl -s -X GET $BASE_URL/quiz/1"
run_curl "$GET_QUIZ_CMD"

# 3. Submit a Correct Answer
SUBMIT_CORRECT_CMD="curl -s -X POST $BASE_URL/quiz/1/answer/user1 \
  -H \"Content-Type: application/json\" \
  -d '{
    \"question_id\": \"q1\",
    \"selected_option\": 1
  }'"

run_curl "$SUBMIT_CORRECT_CMD"

# 4. Submit an Incorrect Answer
SUBMIT_INCORRECT_CMD="curl -s -X POST $BASE_URL/quiz/1/answer/user1 \
  -H \"Content-Type: application/json\" \
  -d '{
    \"question_id\": \"q2\",
    \"selected_option\": 1
  }'"

run_curl "$SUBMIT_INCORRECT_CMD"

# 5. Get Results for User1
GET_RESULTS_CMD="curl -s -X GET $BASE_URL/quiz/1/results/user1"
run_curl "$GET_RESULTS_CMD"

# 6. Try to Get a Non-existent Quiz
GET_NONEXISTENT_QUIZ_CMD="curl -s -X GET $BASE_URL/quiz/999"
run_curl "$GET_NONEXISTENT_QUIZ_CMD"

# 7. Submit an Answer with Invalid JSON
SUBMIT_INVALID_CMD="curl -s -X POST $BASE_URL/quiz/1/answer/user2 \
  -H \"Content-Type: application/json\" \
  -d 'invalid json'"

run_curl "$SUBMIT_INVALID_CMD"

# 8. Get Results for a User Who Hasn't Taken the Quiz
GET_NONEXISTENT_RESULTS_CMD="curl -s -X GET $BASE_URL/quiz/1/results/user3"
run_curl "$GET_NONEXISTENT_RESULTS_CMD"

echo "Test completed"