package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (optional - will use system env vars if not found)
	_ = godotenv.Load()
	
	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	hub := newHub(db)
	go hub.run()

	setupRoutes(hub)

	log.Println("Chat server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}