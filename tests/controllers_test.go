package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"quiz-app/internal/controllers"
	"quiz-app/internal/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of the Storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateQuiz(quiz *models.Quiz) error {
	args := m.Called(quiz)
	return args.Error(0)
}

func (m *MockStorage) GetQuiz(id string) (*models.Quiz, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Quiz), args.Error(1)
}

func (m *MockStorage) SubmitAnswer(quizID, userID string, answer *models.Answer) (bool, string, error) {
	args := m.Called(quizID, userID, answer)
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockStorage) GetResults(quizID, userID string) (*models.Result, error) {
	args := m.Called(quizID, userID)
	return args.Get(0).(*models.Result), args.Error(1)
}

func TestCreateQuiz(t *testing.T) {
	t.Run("Successful quiz creation", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)
		quiz := models.Quiz{
			ID:                "1",
			Title:             "Test Quiz",
			IsNegativeMarking: true,
			Penalty:           0.5,
			Questions: []models.Question{
				{
					ID:            "q1",
					Text:          "What is 1+1?",
					Options:       []string{"1", "2", "3", "4"},
					CorrectOption: 1,
					Marks:         2,
				},
			},
		}

		mockStorage.On("CreateQuiz", &quiz).Return(nil)

		body, _ := json.Marshal(quiz)
		req, _ := http.NewRequest("POST", "/quiz", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		controller.CreateQuiz(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		var response map[string]string
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, "Quiz created successfully", response["message"])
		mockStorage.AssertExpectations(t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		req, _ := http.NewRequest("POST", "/quiz", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		controller.CreateQuiz(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request body\n", rr.Body.String())
	})

	t.Run("Storage error", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		quiz := models.Quiz{ID: "1", Title: "Test Quiz"}
		mockStorage.On("CreateQuiz", &quiz).Return(errors.New("storage error"))

		body, _ := json.Marshal(quiz)
		req, _ := http.NewRequest("POST", "/quiz", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		controller.CreateQuiz(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "Failed to create quiz\n", rr.Body.String())
		mockStorage.AssertExpectations(t)
	})
}

func TestGetQuiz(t *testing.T) {

	t.Run("Successful quiz retrieval", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		quiz := &models.Quiz{
			ID:                "1",
			Title:             "Test Quiz",
			IsNegativeMarking: true,
			Penalty:           0.5,
			Questions: []models.Question{
				{
					ID:            "q1",
					Text:          "What is 1+1?",
					Options:       []string{"1", "2", "3", "4"},
					CorrectOption: 1,
					Marks:         2,
				},
			},
		}

		mockStorage.On("GetQuiz", "1").Return(quiz, nil)

		req, _ := http.NewRequest("GET", "/quiz/1", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/quiz/{id}", controller.GetQuiz)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var retrievedQuiz models.Quiz
		json.Unmarshal(rr.Body.Bytes(), &retrievedQuiz)
		assert.Equal(t, quiz.ID, retrievedQuiz.ID)
		assert.Equal(t, quiz.Title, retrievedQuiz.Title)
		assert.Equal(t, quiz.IsNegativeMarking, retrievedQuiz.IsNegativeMarking)
		assert.Equal(t, quiz.Penalty, retrievedQuiz.Penalty)
		assert.Equal(t, 0, retrievedQuiz.Questions[0].CorrectOption)
		assert.Equal(t, 0, retrievedQuiz.Questions[0].Marks)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Quiz not found", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		mockStorage.On("GetQuiz", "2").Return(&models.Quiz{}, errors.New("quiz not found"))

		req, _ := http.NewRequest("GET", "/quiz/2", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/quiz/{id}", controller.GetQuiz)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "Quiz not found\n", rr.Body.String())
		mockStorage.AssertExpectations(t)
	})
}

func TestSubmitAnswer(t *testing.T) {
	t.Run("Correct answer submission", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		answer := models.Answer{QuestionID: "q1", SelectedOption: 1}
		mockStorage.On("SubmitAnswer", "1", "user1", &answer).Return(true, "", nil)

		body, _ := json.Marshal(answer)
		req, _ := http.NewRequest("POST", "/quiz/1/answer/user1", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/quiz/{quizId}/answer/{userId}", controller.SubmitAnswer)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, true, response["is_correct"])
		assert.NotContains(t, response, "correct_answer")
		mockStorage.AssertExpectations(t)
	})

	t.Run("Incorrect answer submission", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		answer := models.Answer{QuestionID: "q1", SelectedOption: 0}
		mockStorage.On("SubmitAnswer", "1", "user1", &answer).Return(false, "2", nil)

		body, _ := json.Marshal(answer)
		req, _ := http.NewRequest("POST", "/quiz/1/answer/user1", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/quiz/{quizId}/answer/{userId}", controller.SubmitAnswer)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response["is_correct"])
		assert.Equal(t, "2", response["correct_answer"])
		mockStorage.AssertExpectations(t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		req, _ := http.NewRequest("POST", "/quiz/1/question/q1/answer", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/quiz/{quizId}/question/{questionId}/answer", controller.SubmitAnswer)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request body\n", rr.Body.String())
	})

	t.Run("Storage error", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		answer := models.Answer{QuestionID: "q1", SelectedOption: 1}
		mockStorage.On("SubmitAnswer", "1", "user1", &answer).Return(false, "", errors.New("storage error"))

		body, _ := json.Marshal(answer)
		req, _ := http.NewRequest("POST", "/quiz/1/answer/user1", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/quiz/{quizId}/answer/{userId}", controller.SubmitAnswer)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "Failed to submit answer\n", rr.Body.String())
		mockStorage.AssertExpectations(t)
	})
}

func TestGetResults(t *testing.T) {

	t.Run("Successful results retrieval", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		result := &models.Result{
			QuizID: "1",
			UserID: "user1",
			Score:  80,
			Answers: map[string]models.Answer{
				"q1": {QuestionID: "q1", SelectedOption: 1, IsCorrect: true},
				"q2": {QuestionID: "q2", SelectedOption: 2, IsCorrect: false},
			},
		}

		mockStorage.On("GetResults", "1", "user1").Return(result, nil)

		req, _ := http.NewRequest("GET", "/quiz/1/results/user1", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/quiz/{quizId}/results/{userId}", controller.GetResults)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var retrievedResult models.Result
		json.Unmarshal(rr.Body.Bytes(), &retrievedResult)
		assert.Equal(t, result.QuizID, retrievedResult.QuizID)
		assert.Equal(t, result.UserID, retrievedResult.UserID)
		assert.Equal(t, result.Score, retrievedResult.Score)
		assert.Equal(t, len(result.Answers), len(retrievedResult.Answers))
		mockStorage.AssertExpectations(t)
	})

	t.Run("Results not found", func(t *testing.T) {
		mockStorage := new(MockStorage)
		controller := controllers.NewQuizController(mockStorage)

		mockStorage.On("GetResults", "1", "user2").Return(&models.Result{}, errors.New("results not found"))

		req, _ := http.NewRequest("GET", "/quiz/1/results/user2", nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/quiz/{quizId}/results/{userId}", controller.GetResults)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "Results not found\n", rr.Body.String())
		mockStorage.AssertExpectations(t)
	})
}
