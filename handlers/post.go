package post

import (
	"forum/db"
	"html/template"
	"log"
	"net/http"
)

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
