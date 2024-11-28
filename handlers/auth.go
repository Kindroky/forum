package auth

import (
	"fmt"
	"forum/db"
	"log"
	"net/http"
)

// Register handle user registration
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if email == "" || username == "" || password == "" {
		http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
		return
	}

	dbConn := db.GetDBConnection()
	var existingUser string
	err := dbConn.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingUser)
	if err == nil {
		http.Error(w, "Cet email est déjà utilisé", http.StatusConflict)
		return
	}

	_, err = dbConn.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, password)
	if err != nil {
		log.Printf("Erreur lors de l'inscription : %v", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Utilisateur créé avec succès !")
}

// Login handle user login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
		return
	}

	dbConn := db.GetDBConnection()
	var dbPassword string
	err := dbConn.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&dbPassword)
	if err != nil || dbPassword != password {
		http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "some-unique-session-id", // Remplacer par une vraie gestion de session
		HttpOnly: true,
	})

	fmt.Fprintf(w, "Connexion réussie !")
}
