package handlers

import (
	"database/sql"
	"errors"
	"forum/db"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func CommentPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("templates/comment.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := 1 // Replace with actual user session logic
	postID := r.FormValue("postID")
	comment := r.FormValue("comment")

	if comment == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	dbConn := db.GetDBConnection()
	err = AddComment(userID, postIDInt, comment, dbConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AddComment(userID, postID int, comment string, db *sql.DB) error {
	// Ensure user is registered
	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user not registered")
	}

	// Ensure post exists
	var postExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&postExists)
	if err != nil {
		return err
	}
	if !postExists {
		return errors.New("post does not exist")
	}

	// Insert comment into the Comments table
	_, err = db.Exec("INSERT INTO comments (user_id, post_id, content, created_at) VALUES (?, ?, ?, datetime('now'))", userID, postID, comment)
	if err != nil {
		log.Printf("Error inserting comment: %v", err)
		return err
	}
	return err
}
