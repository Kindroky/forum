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
		Error(w, r, http.StatusUnauthorized, "You need to log in to add a post.")
		return
	}

	dbConn := db.GetDBConnection()
	var user User
	err = dbConn.QueryRow(`
        SELECT id, username, LP, session_id 
        FROM users WHERE session_id = ?`, cookie.Value).Scan(&user.ID, &user.Username, &user.LP, &user.SessionID)
	if err != nil {
		log.Printf("Error retrieving user from session: %v", err)
		Error(w, r, http.StatusUnauthorized, "Invalid session. Please log in again.")
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/addpost.html")
		if err != nil {
			log.Printf("Error parsing addpost template: %v", err)
			Error(w, r, http.StatusInternalServerError, "Error loading the Add Post page.")
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
			Error(w, r, http.StatusBadRequest, "All fields are required.")
			return
		}

		if len(categories) > 2 {
			Error(w, r, http.StatusBadRequest, "You can only select up to 2 categories.")
			return
		}

		categoriesStr := strings.Join(categories, ",")

		err := db.CreatePost(title, content, categoriesStr, user.ID)
		if err != nil {
			log.Printf("Error creating post: %v", err)
			Error(w, r, http.StatusInternalServerError, "An error occurred while creating your post.")
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
