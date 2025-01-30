package handlers

import (
	"forum/models"
	"html/template"
	"log"
	"net/http"
)

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, lp := getSessionUserID(r)
	if userID == 0 {
		Error(w, r, http.StatusUnauthorized, "You must be logged in to view liked posts.")
		return
	}
	// Calculate the current user's rank
	currentUser := models.User{
		ID:   userID,
		LP:   lp,
		Rank: getRank(lp), // Dynamically calculate rank
	}
	// Fetch liked posts
	posts, err := GetLikedPosts(userID)
	if err != nil {
		log.Printf("Error retrieving liked posts: %v", err)
		Error(w, r, http.StatusInternalServerError, "Unable to retrieve liked posts.")
		return
	}
	// Calculate ranks for each post's user
	for i := range posts {
		posts[i].User.Rank = getRank(posts[i].User.LP)
	}
	data := models.HomepageData{
		Authenticated: true,
		User:          currentUser, // Ensure updated rank is passed
		Posts:         posts,
	}
	tmpl, err := template.ParseFiles("templates/likedposts.html")
	if err != nil {
		log.Printf("Error parsing liked posts template: %v", err)
		Error(w, r, http.StatusInternalServerError, "Unable to load liked posts page.")
		return
	}
	tmpl.Execute(w, data)
}
func CreatedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, lp := getSessionUserID(r)
	if userID == 0 {
		Error(w, r, http.StatusUnauthorized, "You must be logged in to view created posts.")
		return
	}
	// Calculate the current user's rank
	currentUser := models.User{
		ID:   userID,
		LP:   lp,
		Rank: getRank(lp), // Dynamically calculate rank
	}
	// Fetch created posts
	posts, err := GetCreatedPosts(userID)
	if err != nil {
		log.Printf("Error retrieving created posts: %v", err)
		Error(w, r, http.StatusInternalServerError, "Unable to retrieve created posts.")
		return
	}
	// Update each post's user's rank
	for i := range posts {
		posts[i].User.Rank = getRank(posts[i].User.LP)
	}
	// Pass current user and updated posts to template
	data := models.HomepageData{
		Authenticated: true,
		User:          currentUser, // Ensure updated rank is passed
		Posts:         posts,
	}
	tmpl, err := template.ParseFiles("templates/createdposts.html")
	if err != nil {
		log.Printf("Error parsing created posts template: %v", err)
		Error(w, r, http.StatusInternalServerError, "Unable to load created posts page.")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		Error(w, r, http.StatusInternalServerError, "An error occurred while rendering the created posts page.")
	}
}
