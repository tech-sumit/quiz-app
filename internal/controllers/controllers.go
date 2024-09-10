package controllers

import (
	"encoding/json"
	"net/http"

	"quiz-app/internal/models"
	"quiz-app/internal/storage"

	"github.com/gorilla/mux"
)

// QuizController handles quiz-related operations
type QuizController struct {
	store storage.Storage
}

// NewQuizController creates a new QuizController
func NewQuizController(store storage.Storage) *QuizController {
	return &QuizController{store: store}
}

// CreateQuiz handles the creation of a new quiz
func (c *QuizController) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	var quiz models.Quiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.store.CreateQuiz(&quiz); err != nil {
		http.Error(w, "Failed to create quiz", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Quiz created successfully"})
}

// GetQuiz retrieves a quiz by ID
func (c *QuizController) GetQuiz(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	quizID := params["id"]

	quiz, err := c.store.GetQuiz(quizID)
	if err != nil {
		http.Error(w, "Quiz not found", http.StatusNotFound)
		return
	}

	// Remove correct_option and marks from questions
	for i := range quiz.Questions {
		quiz.Questions[i].CorrectOption = 0
		quiz.Questions[i].Marks = 0
	}

	json.NewEncoder(w).Encode(quiz)
}

// SubmitAnswer handles the submission of an answer to a quiz question
func (c *QuizController) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	quizID := params["quizId"]
	userID := params["userId"]

	var answer models.Answer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	isCorrect, correctAnswer, err := c.store.SubmitAnswer(quizID, userID, &answer)
	if err != nil {
		http.Error(w, "Failed to submit answer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"is_correct": isCorrect,
	}

	if !isCorrect {
		response["correct_answer"] = correctAnswer
	}

	json.NewEncoder(w).Encode(response)
}

// GetResults retrieves the results of a user's quiz attempt
func (c *QuizController) GetResults(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	quizID := params["quizId"]
	userID := params["userId"]

	result, err := c.store.GetResults(quizID, userID)
	if err != nil {
		http.Error(w, "Results not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}
