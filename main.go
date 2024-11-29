package main

import (
	"html/template"
	"log"
	"net/http"
	"forum/handlers"
	"forum/db"
)

func main() {
	// Initialize the database
	dbConn := db.InitDB("db/forum.db")
	defer dbConn.Close()

	// Create tables if they don't exist
	db.CreateTables()
	log.Println("Database initialized successfully.")

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/addpost", handlers.AddPost)
	mux.HandleFunc("/", handlers.Homepage)

	// Start the server
	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
