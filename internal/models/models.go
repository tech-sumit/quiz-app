package models

// Quiz represents a quiz with multiple questions
type Quiz struct {
	ID                string     `json:"id"`
	Title             string     `json:"title"`
	Questions         []Question `json:"questions"`
	IsNegativeMarking bool       `json:"is_negative_marking"`
	Penalty           float32    `json:"penalty"`
}

// Question represents a single question in a quiz
type Question struct {
	ID            string   `json:"id"`
	Text          string   `json:"text"`
	Options       []string `json:"options"`
	CorrectOption int      `json:"correct_option,omitempty"`
	Marks         int      `json:"marks"`
}

// Answer represents a user's answer to a question
type Answer struct {
	QuestionID     string `json:"question_id"`
	SelectedOption int    `json:"selected_option"`
	IsCorrect      bool   `json:"is_correct"`
}

// Result represents the overall result of a user's quiz attempt
type Result struct {
	QuizID  string            `json:"quiz_id"`
	UserID  string            `json:"user_id"`
	Score   float32           `json:"score"`
	Answers map[string]Answer `json:"answers"`
}
