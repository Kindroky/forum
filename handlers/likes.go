package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/db"
	"net/http"
	"strconv"
)

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, r, http.StatusMethodNotAllowed, "Invalid request method.")
		return
	}

	userID := getSessionUserID(r)
	if userID == 0 {
		Error(w, r, http.StatusUnauthorized, "You must be logged in to like or dislike a post.")
		return
	}

	postID, likeType, err := parseLikeRequest(r)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	dbConn := db.GetDBConnection()

	err = handleLikeDislike(userID, postID, likeType, dbConn)
	if err != nil {
		// Check if the error is user-related (e.g., liking own post)
		if err.Error() == "users cannot like or dislike their own posts" {
			Error(w, r, http.StatusForbidden, "A user can't like their own posts. One has to go the hard way to earn LP, keep it up!")
			return
		}

		// Log and respond with internal server error for unexpected issues
		Error(w, r, http.StatusInternalServerError, "An unexpected error occurred while processing your request.")
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
		_, err = dbConn.Exec(`
			INSERT INTO likes (user_id, post_id, like_type) VALUES (?, ?, ?)`,
			userID, postID, likeType)
		if err != nil {
			return err
		}
		if likeType == 1 {
			return UpdateLikesAndDislikes(postID, 1, 0)
		}
		return UpdateLikesAndDislikes(postID, 0, 1)
	} else if err != nil {
		return err
	}
	if existingLikeType != likeType {
		_, err = dbConn.Exec(`
			UPDATE likes SET like_type = ? WHERE user_id = ? AND post_id = ?`,
			likeType, userID, postID)
		if err != nil {
			return err
		}
		if likeType == 1 {
			return UpdateLikesAndDislikes(postID, 1, -1)
		}
		return UpdateLikesAndDislikes(postID, -1, 1)
	}
	_, err = dbConn.Exec(`
		DELETE FROM likes WHERE user_id = ? AND post_id = ?`,
		userID, postID)
	if err != nil {
		return err
	}
	if likeType == 1 {
		return UpdateLikesAndDislikes(postID, -1, 0)
	}
	return UpdateLikesAndDislikes(postID, 0, -1)
}

func UpdateLikesAndDislikes(postID, likeChange, dislikeChange int) error {
	db := db.GetDBConnection()
	var likeCount, dislikeCount int

	// Fetch current counts
	err := db.QueryRow("SELECT likes_count, dislikes_count FROM posts WHERE id=?", postID).Scan(&likeCount, &dislikeCount)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Update counts
	likeCount += likeChange
	dislikeCount += dislikeChange
	_, err = db.Exec("UPDATE posts SET likes_count = ?, dislikes_count = ? WHERE id=?", likeCount, dislikeCount, postID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
