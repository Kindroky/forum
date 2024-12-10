package main

import (
	"forum/db"
	"forum/handlers"
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	dbConn := db.InitDB("db/forum.db")
	defer dbConn.Close()

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/post/", handlers.CommentPostHandler)
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/addpost", handlers.AddPost)
	mux.HandleFunc("/like", handlers.LikePostHandler)
	mux.HandleFunc("/comment", handlers.CommentPostHandler)
	mux.HandleFunc("/detailpost", handlers.DetailPostHandler)
	mux.HandleFunc("/", handlers.Homepage)

	// Start the server
	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
