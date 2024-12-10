package handlers

import (
	"database/sql"
	"forum/db"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Post struct {
	ID            int
	Title         string
	Content       string
	Author        string
	Category      string
	UserID        string
	CreatedAt     string
	LikesCount    int
	DislikesCount int
	Comments      []Comment
	User          User
}

type Comment struct {
	ID        int
	UserID    int
	Content   string
	Author    string
	CreatedAt string
}

func DetailPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the query parameters
	postIDStr := r.URL.Query().Get("id")
	if postIDStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	// Initialize the post model
	var post Post

	// Retrieve post details from the database
	query := `SELECT posts.id, posts.title, posts.content, users.username, posts.created_at, posts.likes_count, posts.dislikes_count
              FROM posts
              JOIN users ON posts.user_id = users.id
              WHERE posts.id = ?`
	dbConn := db.GetDBConnection()
	err = dbConn.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt, &post.LikesCount, &post.DislikesCount)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Retrieve comments associated with the post
	rows, err := dbConn.Query(`
		SELECT comments.id, comments.user_id, comments.content, users.username, comments.created_at
		FROM comments
		JOIN users ON comments.user_id = users.id
		WHERE comments.post_id = ?
		ORDER BY comments.created_at ASC`, postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.Author, &comment.CreatedAt); err != nil {
			http.Error(w, "Error reading comments", http.StatusInternalServerError)
			return
		}
		post.Comments = append(post.Comments, comment)
	}

	// Render the detail post page
	tmpl, err := template.ParseFiles("templates/detailpost.html")
	if err != nil {
		log.Printf("Error parsing detail post template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
