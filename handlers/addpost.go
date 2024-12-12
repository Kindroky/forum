package handlers

import (
	"forum/db"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No session cookie found.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	dbConn := db.GetDBConnection()
	var user User
	err = dbConn.QueryRow(`
        SELECT id, username, LP, session_id 
        FROM users WHERE session_id = ?`, cookie.Value).Scan(&user.ID, &user.Username, &user.LP, &user.SessionID)
	if err != nil {
		log.Printf("Error retrieving user from session: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/addpost.html")
		if err != nil {
			log.Printf("Error parsing addpost template: %v", err)
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["category"]

		if title == "" || content == "" || len(categories) == 0 {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if len(categories) > 2 {
			http.Error(w, "You can only select up to 2 categories", http.StatusBadRequest)
			return
		}

		categoriesStr := strings.Join(categories, ",")

		err := db.CreatePost(title, content, categoriesStr, user.ID)
		if err != nil {
			log.Printf("Error creating post: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
