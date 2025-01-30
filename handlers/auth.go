package handlers

import (
	"forum/db"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("templates/register.html")
		if err != nil {
			log.Printf("Error loading register template: %v", err)
			Error(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}
		t.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		if email == "" || username == "" || password == "" {
			Error(w, r, http.StatusBadRequest, "All fields are required")
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			Error(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}

		// Insert user into database
		dbConn := db.GetDBConnection()
		_, err = dbConn.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, string(hashedPassword))
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			Error(w, r, http.StatusConflict, "This email is already in use")
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Printf("Error loading login template: %v", err)
			Error(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}
		t.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			log.Println("Login error: Missing email or password")
			Error(w, r, http.StatusBadRequest, "All fields are required")
			return
		}

		// Retrieve user from database
		dbConn := db.GetDBConnection()
		var hashedPassword string
		var userID int
		err := dbConn.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &hashedPassword)
		if err != nil {
			log.Printf("Login error: Email not found or DB error: %v", err)
			Error(w, r, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			log.Printf("Login error: Password mismatch for email=%s: %v", email, err)
			Error(w, r, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		// Generate session ID and set cookie
		sessionID := uuid.New().String()
		_, err = dbConn.Exec("UPDATE users SET session_id = ? WHERE id = ?", sessionID, userID)
		if err != nil {
			log.Printf("Error updating session ID: %v", err)
			Error(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
		})

		log.Printf("User logged in successfully: email=%s", email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
