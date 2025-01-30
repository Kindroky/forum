package main

import (
	"forum/db"
	"forum/handlers"
	"log"
	"net/http"
)

func main() {
	dbConn := db.InitDB("db/forum.db")
	defer dbConn.Close()
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/post/", handlers.CommentFormHandler)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/addpost", handlers.AddPost)
	mux.HandleFunc("/like", handlers.LikePostHandler)
	mux.HandleFunc("/comment", handlers.CommentPostHandler)
	mux.HandleFunc("/detailpost", handlers.DetailPostHandler)
	mux.HandleFunc("/comlike", handlers.ComLikePostHandler)
	mux.HandleFunc("/likedposts", handlers.LikedPostsHandler)
	mux.HandleFunc("/createdposts", handlers.CreatedPostsHandler)
	mux.HandleFunc("/", handlers.Homepage)
	defer db.CloseDB()
	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
