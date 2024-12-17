package handlers

import (
	"forum/db"
	"forum/models"
	"html/template"
	"log"
	"net/http"
)

func PostDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL query
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	// Check session for authentication
	authenticated := false
	var user models.User
	cookie, err := r.Cookie("session_id")
	if err == nil {
		dbConn := db.GetDBConnection()
		err = dbConn.QueryRow(`
			SELECT id, username, 
			       CASE 
			           WHEN LP >= 100 THEN 'Legend' 
			           WHEN LP >= 50 THEN 'Pro' 
			           ELSE 'Novice' 
			       END AS rank, LP, session_id 
			FROM users WHERE session_id = ?`, cookie.Value).Scan(
			&user.ID, &user.Username, &user.Rank, &user.LP, &user.SessionID)
		if err == nil {
			authenticated = true
		} else {
			log.Printf("Error fetching user data: %v", err)
		}
	}
	// Fetch the post details
	dbConn := db.GetDBConnection()
	var post models.Post
	err = dbConn.QueryRow(`
		SELECT 
			posts.id, 
			posts.title, 
			posts.content, 
			posts.category, 
			posts.likes_count, 
			posts.dislikes_count, 
			posts.user_id, 
			users.username, 
			CASE 
				WHEN users.LP >= 100 THEN 'Legend' 
				WHEN users.LP >= 50 THEN 'Pro' 
				ELSE 'Novice' 
			END AS rank, 
			users.LP
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.id = ?`, postID).Scan(
		&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &post.User.Username, &post.User.Rank, &post.User.LP)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	// Fetch the comments for this post
	comments, err := FetchComments(post.ID)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Create data for the template
	data := struct {
		Authenticated bool
		User          models.User
		Post          models.Post
		Comments      []models.Comment
	}{
		Authenticated: authenticated,
		User:          user,
		Post:          post,
		Comments:      comments,
	}
	// Render the post details template
	tmpl, err := template.ParseFiles("templates/postdetails.html")
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
