package handlers

import (
	"forum/db"
	"html/template"
	"log"
	"net/http"
)

type HomepageData struct {
	Authenticated bool
	User          User
	Posts         []Post
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user User

	cookie, err := r.Cookie("session_id")
	if err == nil {
		dbConn := db.GetDBConnection()
		err = dbConn.QueryRow(`
            SELECT id, username, LP, session_id 
            FROM users WHERE session_id = ?`, cookie.Value).Scan(&user.ID, &user.Username, &user.LP, &user.SessionID)
		if err == nil {
			authenticated = true
		} else {
			log.Printf("Error fetching user data: %v", err)
		}
	}

	categories := r.URL.Query()["category"]

	var posts []Post
	if len(categories) > 0 {
		posts, err = GetPostsByCategory(categories)
	} else {
		posts, err = GetPosts()
	}

	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := HomepageData{
		Authenticated: authenticated,
		User:          user,
		Posts:         posts,
	}

	tmpl, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
