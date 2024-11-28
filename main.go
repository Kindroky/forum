package main

import (
    "html/template"
    "forum/handlers"
    "forum/db"
    "log"
    "net/http"
)

func main() {
    // Initialisation de la base de données
    database := db.InitDB("db/forum.db")
    defer database.Close()

    	// Création des tables
	db.CreateTables()

    log.Println("Base de données initialisée avec succès.")

    // Routes
    mux := http.NewServeMux()
    mux.HandleFunc("/register", handlers.Register)
    mux.HandleFunc("/login", handlers.Login)
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("templates/home.html")
        if err != nil {
            http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
            return
        }
        posts, err := handlers.GetPosts()
        if err!=nil {
            http.Error(w, "Erreur lors de la recuperation des posts", http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, posts)
    })
    

    // Démarrage du serveur
    log.Println("Le serveur est lancé sur http://localhost:8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalf("Erreur lors du démarrage du serveur : %v", err)
    }
}
