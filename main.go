package main

import (
	"forum/db"
	"forum/handlers"
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	database := db.InitDB("db/forum.db")
	defer database.Close()

	// Create tables if they don't exist
	db.CreateTables()
	log.Println("Database initialized successfully.")

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/addpost", handlers.AddPost)
	mux.HandleFunc("/", handlers.Homepage) // Use the dedicated Homepage handler

	// Start the server
	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
