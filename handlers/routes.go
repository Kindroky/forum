package handlers

import (
	"encoding/json"
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
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles(`templates/register.html`)
		if err != nil {
			http.Error(w, "internal serveur error", http.StatusInternalServerError)
		}
		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
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

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
}
// Login handle user login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
	} else if r.Method == http.MethodGet {
		if r.Method == http.MethodGet {
			t, err := template.ParseFiles(`templates/login.html`)
			if err != nil {
				http.Error(w, "internal serveur error", http.StatusInternalServerError)
			}
			t.Execute(w, nil)
		}
	}
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
		category := r.FormValue("category")

		if title == "" || content == "" || category == "" {
			http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
			return
		}

		validCategories := map[string]bool{"items": true, "champs": true, "meta": true, "events": true}
		if !validCategories[category] {
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}

		err := db.CreatePost(title, content, category, 1) // Assuming user_id = 1
		if err != nil {
			log.Printf("Erreur lors de l'ajout du post : %v", err)
			http.Error(w, "Erreur interne", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Homepage displays posts with optional category filtering
func Homepage(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var posts []db.Post
	var err error

	if category == "" {
		posts, err = db.GetPosts()
	} else {
		posts, err = db.GetPostsByCategory(category)
	}

	if err != nil {
		log.Printf("Erreur lors de la récupération des posts : %v", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		log.Printf("Erreur lors du chargement du template : %v", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Posts []db.Post
	}{Posts: posts})
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template : %v", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
	}
}

func GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	validCategories := map[string]bool{"items": true, "champs": true, "meta": true, "events": true}

	if !validCategories[category] {
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}

	// Use the database connection from the db package
	dbConn := db.GetDBConnection()

	rows, err := dbConn.Query("SELECT id, title, content, category, user_id FROM posts WHERE category = ? ORDER BY id DESC", category)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Serialize the posts and return them as JSON
	posts := []db.Post{}
	for rows.Next() {
		var post db.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID); err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}
