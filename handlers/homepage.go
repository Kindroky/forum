package handlers

import (
	"forum/db"
	"html/template"
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
		Error(w, r, http.StatusInternalServerError, "An error occurred while fetching posts.")
		return
	}

	data := HomepageData{
		Authenticated: authenticated,
		User:          user,
		Posts:         posts,
	}

	tmpl, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		Error(w, r, http.StatusInternalServerError, "An error occurred while loading the homepage.")
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, "An error occurred while rendering the homepage.")
	}
}
