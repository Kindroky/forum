package handlers

import (
	"fmt"
	"forum/db"
	"html/template"
	"log"
	"net/http"
)

type Post struct {
	ID      int
	Title   string
	Content string
	UserID  int
}

// Register handles user registration
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

// Login handles user login
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
		Value:    "some-unique-session-id", // Replace with actual session management logic
		HttpOnly: true,
	})

	fmt.Fprintf(w, "Connexion réussie !")
}

// AddPost allows users to create a post
func AddPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/addpost.html")
		if err != nil {
			http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		if title == "" || content == "" {
			http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
			return
		}

		err := db.CreatePost(title, content, 1) // Assuming user_id = 1
		if err != nil {
			log.Printf("Erreur lors de l'ajout du post : %v", err)
			http.Error(w, "Erreur interne", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Homepage displays all posts on the homepage
func Homepage(w http.ResponseWriter, r *http.Request) {
	posts, err := db.GetPosts()
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		log.Printf("Error loading homepage template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Posts []db.Post
	}{Posts: posts})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
