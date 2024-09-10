package tests

import (
	"testing"

	"quiz-app/internal/models"
	"quiz-app/internal/storage"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_CreateAndGetQuiz(t *testing.T) {
	store := storage.NewMemoryStorage()

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

	err := store.CreateQuiz(quiz)
	assert.NoError(t, err)

	retrievedQuiz, err := store.GetQuiz("1")
	assert.NoError(t, err)
	assert.Equal(t, quiz.ID, retrievedQuiz.ID)
	assert.Equal(t, quiz.Title, retrievedQuiz.Title)
	assert.Equal(t, quiz.IsNegativeMarking, retrievedQuiz.IsNegativeMarking)
	assert.Equal(t, quiz.Penalty, retrievedQuiz.Penalty)
	assert.Equal(t, len(quiz.Questions), len(retrievedQuiz.Questions))
	assert.Equal(t, quiz.Questions[0].ID, retrievedQuiz.Questions[0].ID)
	assert.Equal(t, quiz.Questions[0].Marks, retrievedQuiz.Questions[0].Marks)
}

func TestMemoryStorage_SubmitAnswerAndGetResults(t *testing.T) {
	store := storage.NewMemoryStorage()

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
			{
				ID:            "q2",
				Text:          "What is 2+2?",
				Options:       []string{"2", "3", "4", "5"},
				CorrectOption: 2,
				Marks:         3,
			},
		},
	}

	err := store.CreateQuiz(quiz)
	assert.NoError(t, err)

	// Test correct answer
	isCorrect, correctAnswer, err := store.SubmitAnswer("1", "user1", &models.Answer{QuestionID: "q1", SelectedOption: 1})
	assert.NoError(t, err)
	assert.True(t, isCorrect)
	assert.Empty(t, correctAnswer)

	// Test incorrect answer
	isCorrect, correctAnswer, err = store.SubmitAnswer("1", "user1", &models.Answer{QuestionID: "q2", SelectedOption: 1})
	assert.NoError(t, err)
	assert.False(t, isCorrect)
	assert.Equal(t, "4", correctAnswer)

	// Get results
	result, err := store.GetResults("1", "user1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.QuizID)
	assert.Equal(t, "user1", result.UserID)
	assert.Equal(t, float32(1.5), result.Score) // 2 (correct) - 0.5 (penalty)
	assert.Equal(t, 2, len(result.Answers))
	assert.True(t, result.Answers["q1"].IsCorrect)
	assert.False(t, result.Answers["q2"].IsCorrect)
}

func TestMemoryStorage_GetResultsNotFound(t *testing.T) {
	store := storage.NewMemoryStorage()

	// Test for non-existent quiz
	_, err := store.GetResults("nonexistent", "user1")
	assert.Error(t, err)
	assert.Equal(t, "no results found for this quiz", err.Error())

	// Create a quiz but don't submit any answers
	quiz := &models.Quiz{ID: "1", Title: "Test Quiz", Questions: []models.Question{
		{
			ID:            "1",
			Text:          "abc",
			Options:       []string{"A", "B", "C", "D"},
			CorrectOption: 1,
			Marks:         1,
		},
	}}
	store.CreateQuiz(quiz)
	store.SubmitAnswer("1", "userX", &models.Answer{
		QuestionID:     "1",
		SelectedOption: 1,
		IsCorrect:      true,
	})

	// Test for non-existent user
	_, err = store.GetResults("1", "nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "no results found for this user", err.Error())
}

func TestMemoryStorage_SubmitAnswerQuizNotFound(t *testing.T) {
	store := storage.NewMemoryStorage()

	_, _, err := store.SubmitAnswer("nonexistent", "user1", &models.Answer{QuestionID: "q1", SelectedOption: 0})
	assert.Error(t, err)
	assert.Equal(t, "quiz not found", err.Error())
}

func TestMemoryStorage_SubmitAnswerQuestionNotFound(t *testing.T) {
	store := storage.NewMemoryStorage()

	quiz := &models.Quiz{
		ID:    "1",
		Title: "Test Quiz",
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

	store.CreateQuiz(quiz)

	_, _, err := store.SubmitAnswer("1", "user1", &models.Answer{QuestionID: "nonexistent", SelectedOption: 0})
	assert.Error(t, err)
	assert.Equal(t, "question not found", err.Error())
}
