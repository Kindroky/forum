package handlers

import (
	"forum/db"
	"net/http"
	"strconv"
)

func CommentPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := getSessionUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
	_, err = dbConn.Exec(`
		INSERT INTO comments (post_id, user_id, content, created_at)
		VALUES (?, ?, ?, datetime('now'))`, postIDInt, userID, comment)
	if err != nil {
		http.Error(w, "Error adding comment", http.StatusInternalServerError)
		return
	}

	err = db.UpdateCommentsCount(postIDInt)
	if err != nil {
		http.Error(w, "Error updating comments count", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/detailpost?id="+postID, http.StatusSeeOther)
}

func FetchComments(postID int) ([]Comment, error) {
	dbConn := db.GetDBConnection()

	// Query to fetch comments joined with user data
	rows, err := dbConn.Query(`
		SELECT comments.id, comments.user_id, comments.content, users.username, comments.created_at
		FROM comments
		JOIN users ON comments.user_id = users.id
		WHERE comments.post_id = ?
		ORDER BY comments.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to store all comments
	var comments []Comment
	for rows.Next() {
		var comment Comment
		// Scan the row into the comment struct
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.Author, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	// Check for errors after iterating through rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
