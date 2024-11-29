package handlers

import (
	"database/sql"
	"errors"
	"forum/db"
	"net/http"
	"strconv"
)

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := 1 // Replace with actual user session logic
	postID := r.FormValue("postID")
	likeType := r.FormValue("likeType")

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	likeTypeInt, err := strconv.Atoi(likeType)
	if err != nil || (likeTypeInt != 1 && likeTypeInt != -1) {
		http.Error(w, "Invalid like type", http.StatusBadRequest)
		return
	}

	dbConn := db.GetDBConnection()
	err = HandleLikes(userID, postIDInt, likeTypeInt, dbConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleLikes(userID, postID, likeType int, db *sql.DB) error {
	// Ensure user is registered
	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user not registered")
	}

	// Check if user has already liked/disliked this post
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE UserID = ? AND PostID = ?)", userID, postID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user has already liked or disliked this post")
	}

	// Insert like/dislike into the Likes table
	_, err = db.Exec("INSERT INTO likes (UserID, PostID, LikeType) VALUES (?, ?, ?)", userID, postID, likeType)
	if err != nil {
		return err
	}

	// Update LP for the post owner
	var postOwnerID int
	err = db.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&postOwnerID)
	if err != nil {
		return err
	}

	lpChange := 10
	if likeType == -1 {
		lpChange = -10
	}
	_, err = db.Exec("UPDATE users SET LP = LP + ? WHERE id = ?", lpChange, postOwnerID)
	return err
}
