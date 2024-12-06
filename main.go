package main

import (
	"forum/db"
	"forum/handlers"
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	//comment for checkpoint commit
	dbConn := db.InitDB("db/forum.db")
	defer dbConn.Close()

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/addpost", handlers.AddPost)
	mux.HandleFunc("/", handlers.Homepage)

	// Start the server
	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
