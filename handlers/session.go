package handlers

import (
	"database/sql"
	"forum/db"
	"log"
	"net/http"
)

func getRank(lp int) string {
	switch {
	case lp < 10:
		return "Iron"
	case lp < 20:
		return "Bronze"
	case lp < 40:
		return "Silver"
	case lp < 80:
		return "Gold"
	case lp < 140:
		return "Platinum"
	case lp < 200:
		return "Emerald"
	case lp < 280:
		return "Diamond"
	case lp < 360:
		return "Master"
	case lp < 500:
		return "Grandmaster"
	default:
		return "Challenger"
	}
}
func getSessionUserID(r *http.Request) (int, string, int) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("Session cookie error: %v", err)
		return 0, "", 0
	}
	sessionID := cookie.Value
	if sessionID == "" {
		log.Println("Empty session ID in cookie")
		return 0, "", 0
	}
	dbConn := db.GetDBConnection()
	var userID int
	var lp int
	err = dbConn.QueryRow("SELECT id, LP FROM users WHERE session_id = ?", sessionID).Scan(&userID, &lp)
	if err == sql.ErrNoRows {
		log.Println("Session ID not found in database")
		return 0, "", 0
	} else if err != nil {
		log.Printf("Database error validating session ID: %v", err)
		return 0, "", 0
	}
	rank := getRank(lp) // Calculate the rank based on LP
	return userID, rank, lp
}
