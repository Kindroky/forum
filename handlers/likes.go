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

	userID := getSessionUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, likeType, err := parseLikeRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConn := db.GetDBConnection()
	err = handleLikeDislike(userID, postID, likeType, dbConn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func parseLikeRequest(r *http.Request) (int, int, error) {
	postID, err := strconv.Atoi(r.FormValue("postID"))
	if err != nil {
		return 0, 0, errors.New("invalid post ID")
	}
	likeType, err := strconv.Atoi(r.FormValue("likeType"))
	if err != nil || (likeType != 1 && likeType != -1) {
		return 0, 0, errors.New("invalid like type")
	}
	return postID, likeType, nil
}

func handleLikeDislike(userID, postID, likeType int, dbConn *sql.DB) error {
	var postOwnerID int
	err := dbConn.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&postOwnerID)
	if err != nil {
		return err
	}
	if postOwnerID == userID {
		return errors.New("users cannot like or dislike their own posts")
	}

	var existingLikeType int
	err = dbConn.QueryRow(`
		SELECT like_type FROM likes WHERE user_id = ? AND post_id = ?`,
		userID, postID).Scan(&existingLikeType)

	if err == sql.ErrNoRows {
		// Add a new like/dislike
		_, err = dbConn.Exec(`
			INSERT INTO likes (user_id, post_id, like_type) VALUES (?, ?, ?)`,
			userID, postID, likeType)
		if err != nil {
			return err
		}
		if likeType == 1 {
			return db.UpdateLikesAndDislikes(postID, 1, 0)
		}
		return db.UpdateLikesAndDislikes(postID, 0, 1)
	} else if err != nil {
		return err
	}

	// Reaction exists; handle updates
	if existingLikeType != likeType {
		// Change the type of reaction (like -> dislike or dislike -> like)
		_, err = dbConn.Exec(`
			UPDATE likes SET like_type = ? WHERE user_id = ? AND post_id = ?`,
			likeType, userID, postID)
		if err != nil {
			return err
		}
		if likeType == 1 {
			return db.UpdateLikesAndDislikes(postID, 1, -1)
		}
		return db.UpdateLikesAndDislikes(postID, -1, 1)
	}

	// Reaction exists and matches; remove it
	_, err = dbConn.Exec(`
		DELETE FROM likes WHERE user_id = ? AND post_id = ?`,
		userID, postID)
	if err != nil {
		return err
	}
	if likeType == 1 {
		return db.UpdateLikesAndDislikes(postID, -1, 0)
	}
	return db.UpdateLikesAndDislikes(postID, 0, -1)
}

func getSessionUserID(r *http.Request) int {
	// Replace with actual session logic
	return 1
}
