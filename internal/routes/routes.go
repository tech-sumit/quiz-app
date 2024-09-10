package routes

import (
	"quiz-app/internal/controllers"
	"quiz-app/internal/storage"

	"github.com/gorilla/mux"
)

// SetupRoutes configures and returns the application router
func SetupRoutes(store storage.Storage) *mux.Router {
	r := mux.NewRouter()
	c := controllers.NewQuizController(store)

	r.HandleFunc("/quiz", c.CreateQuiz).Methods("POST")
	r.HandleFunc("/quiz/{id}", c.GetQuiz).Methods("GET")
	r.HandleFunc("/quiz/{quizId}/answer/{userId}", c.SubmitAnswer).Methods("POST")
	r.HandleFunc("/quiz/{quizId}/results/{userId}", c.GetResults).Methods("GET")

	return r
}
