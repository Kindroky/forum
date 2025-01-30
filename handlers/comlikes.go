package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/db"
	"net/http"
	"strconv"
)

func ComLikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, r, http.StatusMethodNotAllowed, "Invalid request method.")
		return
	}

	userID, _, _ := getSessionUserID(r)
	if userID == 0 {
		Error(w, r, http.StatusUnauthorized, "You must be logged in to like or dislike a comment.")
		return
	}

	commentID, comlikeType, err := parseComLikeRequest(r)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	dbConn := db.GetDBConnection()

	err = handleComLikeDislike(userID, commentID, comlikeType, dbConn)
	if err != nil {
		// Check if the error is user-related (e.g., liking own post)
		if err.Error() == "users cannot like or dislike their own comments" {
			Error(w, r, http.StatusForbidden, "A user can't like their own comments. One has to go the hard way to earn LP, keep it up!")
			return
		}

		// Log and respond with internal server error for unexpected issues
		Error(w, r, http.StatusInternalServerError, "An unexpected error occurred while processing your request.")
		return
	}

	var postID int
	err = dbConn.QueryRow("SELECT post_id FROM comments WHERE id = ?", commentID).Scan(&postID)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, "Unable to find the post for this comment.")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/detailpost?id=%d", postID), http.StatusSeeOther)
}

func parseComLikeRequest(r *http.Request) (int, int, error) {
	commentID, err := strconv.Atoi(r.FormValue("commentID"))
	if err != nil {
		return 0, 0, errors.New("invalid comment ID")
	}

	comlikeType, err := strconv.Atoi(r.FormValue("comlikeType"))
	if err != nil || (comlikeType != 1 && comlikeType != -1) {
		return 0, 0, errors.New("invalid like type")
	}

	return commentID, comlikeType, nil
}

func handleComLikeDislike(userID, commentID, comlikeType int, dbConn *sql.DB) error {
	var commentOwnerID int
	err := dbConn.QueryRow("SELECT user_id FROM comments WHERE id = ?", commentID).Scan(&commentOwnerID)
	if err != nil {
		return err
	}
	if commentOwnerID == userID {
		return errors.New("users cannot like or dislike their own comments")
	}
	var existingcomLikeType int
	err = dbConn.QueryRow(`
		SELECT comlike_type FROM comlikes WHERE user_id = ? AND comment_id = ?`,
		userID, commentID).Scan(&existingcomLikeType)
	if err == sql.ErrNoRows {
		_, err = dbConn.Exec(`
			INSERT INTO comlikes (user_id, comment_id, comlike_type) VALUES (?, ?, ?)`,
			userID, commentID, comlikeType)
		if err != nil {
			return err
		}
		if comlikeType == 1 {
			return UpdateComLikesAndDislikes(commentID, 1, 0)
		}
		return UpdateComLikesAndDislikes(commentID, 0, 1)
	} else if err != nil {
		return err
	}
	if existingcomLikeType != comlikeType {
		_, err = dbConn.Exec(`
			UPDATE comlikes SET comlike_type = ? WHERE user_id = ? AND comment_id = ?`,
			comlikeType, userID, commentID)
		if err != nil {
			return err
		}
		if comlikeType == 1 {
			return UpdateComLikesAndDislikes(commentID, 1, -1)
		}
		return UpdateComLikesAndDislikes(commentID, -1, 1)
	}
	_, err = dbConn.Exec(`
		DELETE FROM comlikes WHERE user_id = ? AND comment_id = ?`,
		userID, commentID)
	if err != nil {
		return err
	}
	if comlikeType == 1 {
		return UpdateComLikesAndDislikes(commentID, -1, 0)
	}
	return UpdateComLikesAndDislikes(commentID, 0, -1)
}

func UpdateComLikesAndDislikes(commentID, comlikeChange, comdislikeChange int) error {
	db := db.GetDBConnection()

	// Fetch current counts
	var comLikesCount, comDislikesCount int
	err := db.QueryRow("SELECT comlikes_count, comdislikes_count FROM comments WHERE id=?", commentID).
		Scan(&comLikesCount, &comDislikesCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("comment not found")
		}
		fmt.Println("Error fetching comment counts:", err)
		return err
	}

	// Increment the counts
	comLikesCount += comlikeChange
	comDislikesCount += comdislikeChange

	// Update the counts
	_, err = db.Exec("UPDATE comments SET comlikes_count = ?, comdislikes_count = ? WHERE id=?", comLikesCount, comDislikesCount, commentID)
	if err != nil {
		fmt.Println("Error updating comment counts:", err)
		return err
	}

	return nil
}
