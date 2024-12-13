package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	type StrucErreur struct {
		CodeErr    int
		MessageErr string
	}

	// Populate the error structure
	Erreur := StrucErreur{
		CodeErr:    statusCode,
		MessageErr: message,
	}

	// Parse the error template
	t, err := template.ParseFiles("templates/error.html")
	if err != nil {
		log.Printf("Error loading error template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Write the status code and render the error template
	w.WriteHeader(statusCode) // Ensure the status code is written
	err = t.Execute(w, Erreur)
	if err != nil {
		log.Printf("Error rendering error template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
