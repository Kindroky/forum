package main

import (
	"html/template"
	"log"
	"net/http"
	"forum/handlers"
	"forum/db"
)

func main() {
	// Initialisation de la base de données
	database := db.InitDB("db/forum.db")
	defer database.Close()

	// Création des tables
	db.CreateTables()

	log.Println("Base de données initialisée avec succès.")

	// Initialisation du routeur
	mux := http.NewServeMux()

	// Définition des routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/home.html")
		if err != nil {
			http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
			return
		}
		posts, err := handlers.GetPosts()
		if err != nil {
			log.Printf("Erreur lors de la récupération des posts : %v", err)
			http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, posts)
	})

	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/addpost", handlers.AddPost)

	// Démarrage du serveur
	log.Println("Le serveur est lancé sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur : %v", err)
	}
}
