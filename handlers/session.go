package handlers

import (
	"database/sql"
	"forum/db"
	"log"
	"net/http"
)

func getSessionUserID(r *http.Request) int {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("Session cookie error: %v", err)
		return 0
	}

	sessionID := cookie.Value
	if sessionID == "" {
		log.Println("Empty session ID in cookie")
		return 0
	}

	dbConn := db.GetDBConnection()

	var userID int
	err = dbConn.QueryRow("SELECT id FROM users WHERE session_id = ?", sessionID).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Println("Session ID not found in database")
		return 0
	} else if err != nil {
		log.Printf("Database error validating session ID: %v", err)
		return 0
	}

	return userID
}
