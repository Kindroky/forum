package handlers

import (
	"forum/db"
	"forum/models"
	"html/template"
	"net/http"
	"strconv"
)

func CommentPostHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve user session ID (to confirm the user is logged in)
	userID := getSessionUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract post ID and comment from the form
	postID := r.FormValue("ID")
	comment := r.FormValue("comment")

	if comment == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	// Convert postID to an integer
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Insert the comment into the database
	dbConn := db.GetDBConnection()
	_, err = dbConn.Exec(`
		INSERT INTO comments (post_id, user_id, content, created_at)
		VALUES (?, ?, ?, datetime('now'))`, postIDInt, userID, comment)
	if err != nil {
		http.Error(w, "Error adding comment", http.StatusInternalServerError)
		return
	}

	// Update the comment count
	err = db.UpdateCommentsCount(postIDInt)
	if err != nil {
		http.Error(w, "Error updating comments count", http.StatusInternalServerError)
		return
	}

	// Redirect back to the detail post page
	http.Redirect(w, r, "/detailpost?id="+postID, http.StatusSeeOther)
}

func FetchComments(postID int) ([]models.Comment, error) {
	dbConn := db.GetDBConnection()
	rows, err := dbConn.Query(`
		SELECT comments.id, comments.user_id, comments.content, users.username, comments.created_at
		FROM comments
		JOIN users ON comments.user_id = users.id
		WHERE comments.post_id = ?
		ORDER BY comments.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.Author, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
func CommentFormHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("ID")
	post, err := GetPostById(postID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	tmpl, err := template.ParseFiles("templates/comment.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	tmpl.Execute(w, post)
}
