package storage

import (
	"errors"
	"sync"

	"quiz-app/internal/models"
)

type Storage interface {
	CreateQuiz(quiz *models.Quiz) error
	GetQuiz(id string) (*models.Quiz, error)
	SubmitAnswer(quizID, userID string, answer *models.Answer) (bool, string, error)
	GetResults(quizID, userID string) (*models.Result, error)
}

type MemoryStorage struct {
	quizzes map[string]models.Quiz
	results map[string]map[string]models.Result // map[quizID]map[userID]Result
	mu      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		quizzes: make(map[string]models.Quiz),
		results: make(map[string]map[string]models.Result),
	}
}

func (m *MemoryStorage) CreateQuiz(quiz *models.Quiz) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.quizzes[quiz.ID] = *quiz
	return nil
}

func (m *MemoryStorage) GetQuiz(id string) (*models.Quiz, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	quiz, exists := m.quizzes[id]
	if !exists {
		return nil, errors.New("quiz not found")
	}
	return &quiz, nil
}

func (m *MemoryStorage) SubmitAnswer(quizID, userID string, answer *models.Answer) (bool, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	quiz, exists := m.quizzes[quizID]
	if !exists {
		return false, "", errors.New("quiz not found")
	}

	var question models.Question
	for _, q := range quiz.Questions {
		if q.ID == answer.QuestionID {
			question = q
			break
		}
	}

	if question.ID == "" {
		return false, "", errors.New("question not found")
	}

	isCorrect := answer.SelectedOption == question.CorrectOption
	answer.IsCorrect = isCorrect

	// Initialize results for this quiz if not exist
	if m.results[quizID] == nil {
		m.results[quizID] = make(map[string]models.Result)
	}

	// Get or initialize user's result
	result, exists := m.results[quizID][userID]
	if !exists {
		result = models.Result{
			QuizID:  quizID,
			UserID:  userID,
			Score:   0,
			Answers: make(map[string]models.Answer),
		}
	}

	// Update score
	if isCorrect {
		result.Score += float32(question.Marks)
	} else if quiz.IsNegativeMarking {
		result.Score -= quiz.Penalty
	}

	// Store the answer
	result.Answers[answer.QuestionID] = *answer

	// Update the result in storage
	m.results[quizID][userID] = result

	if isCorrect {
		return true, "", nil
	}

	// Return the correct answer option
	if question.CorrectOption >= 0 && question.CorrectOption < len(question.Options) {
		return false, question.Options[question.CorrectOption], nil
	}
	return false, "", errors.New("invalid correct option")
}

func (m *MemoryStorage) GetResults(quizID, userID string) (*models.Result, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	quizResults, exists := m.results[quizID]
	if !exists {
		return nil, errors.New("no results found for this quiz")
	}

	result, exists := quizResults[userID]
	if !exists {
		return nil, errors.New("no results found for this user")
	}

	return &result, nil
}
