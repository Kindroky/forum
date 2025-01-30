package handlers

import (
	"database/sql"
	"forum/db"
	"forum/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func DetailPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the query parameters
	postIDStr := r.URL.Query().Get("id")
	if postIDStr == "" {
		Error(w, r, http.StatusBadRequest, "Post ID is required.")
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		Error(w, r, http.StatusBadRequest, "Invalid Post ID.")
		return
	}

	// Initialize the post model
	var post models.Post
	var user models.User

	// Retrieve post details from the database
	query := `SELECT posts.id, posts.title, posts.content, posts.created_at, posts.likes_count, posts.dislikes_count, users.username, users.LP
              FROM posts
              JOIN users ON posts.user_id = users.id
              WHERE posts.id = ?`
	dbConn := db.GetDBConnection()
	err = dbConn.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.LikesCount, &post.DislikesCount, &user.Username, &user.LP)
	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, r, http.StatusNotFound, "Post not found.")
		} else {
			Error(w, r, http.StatusInternalServerError, "An error occurred while fetching the post.")
		}
		return
	}

	// Retrieve comments associated with the post
	rows, err := dbConn.Query(`
		SELECT comments.id, comments.user_id, comments.content, users.username, comments.created_at, comments.comlikes_count, comments.comdislikes_count
		FROM comments
		JOIN users ON comments.user_id = users.id
		WHERE comments.post_id = ?
		ORDER BY comments.created_at ASC`, postID)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, "Failed to fetch comments.")
		return
	}

	// Assign the rank dynamically
	user.Rank = getRank(user.LP)
	post.User = user // Attach the user to the post

	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.Content,
			&comment.Author,
			&comment.CreatedAt,
			&comment.ComLikesCount,
			&comment.ComDislikesCount,
		); err != nil {
			Error(w, r, http.StatusInternalServerError, "An error occurred while reading comments.")
			return
		}
		post.Comments = append(post.Comments, comment)
	}

	if err = rows.Err(); err != nil {
		Error(w, r, http.StatusInternalServerError, "An error occurred while processing comments.")
		return
	}

	// Parse and render the template
	tmpl, err := template.ParseFiles("templates/detailpost.html")
	if err != nil {
		log.Printf("Error parsing detail post template: %v", err)
		http.Error(w, "An error occurred while loading the post page.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "An error occurred while rendering the post page.", http.StatusInternalServerError)
		return
	}
}
