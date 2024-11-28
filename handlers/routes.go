package handlers

import (
	"html/template"
	"forum/db"
	"fmt"
	"log"
	"net/http"
)

type Post struct {
	ID      int
	Title   string
	Content string
	UserID  int
}

// GetPosts récupère tous les posts
func GetPosts() ([]Post, error) {
	dbConn := db.GetDBConnection()
	rows, err := dbConn.Query("SELECT * FROM posts")
	if err != nil {
		log.Printf("Erreur lors de la récupération des posts : %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
			log.Printf("Erreur lors du scan des posts : %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

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

// ForumPage affiche tous les posts
/*func ForumPage(w http.ResponseWriter, r *http.Request) {
	posts, err := GetPosts()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
		return
	}
	
}*/

// AddPost permet aux utilisateurs d'ajouter un post
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

        dbConn := db.GetDBConnection()
        _, err := dbConn.Exec("INSERT INTO posts (title, content, user_id) VALUES (?, ?, ?)", title, content, 1)
        if err != nil {
            log.Printf("Erreur lors de l'ajout du post : %v", err)
            http.Error(w, "Erreur interne", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}
