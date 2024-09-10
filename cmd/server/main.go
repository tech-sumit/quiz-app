package main

import (
	"fmt"
	"log"
	"net/http"

	"quiz-app/internal/routes"
	"quiz-app/internal/storage"
)

func main() {
	// Initialize storage
	store := storage.NewMemoryStorage()

	// Initialize router
	router := routes.SetupRoutes(store)

	// Start server
	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
