package handlers

import (
	"forum/db"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	LP        int
	SessionID string
}

func GetPostsByCategory(categories []string) ([]Post, error) {
	dbConn := db.GetDBConnection()

	// Build the query dynamically based on the number of categories
	query := `
		SELECT posts.id, posts.title, posts.content, posts.category, posts.likes_count, posts.dislikes_count, posts.user_id, users.username, users.LP 
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE `

	conditions := make([]string, len(categories))
	args := []any{}

	for i, category := range categories {
		conditions[i] = "category LIKE ?"
		args = append(args, "%"+category+"%") // Add wildcard to match partial values
	}

	query += strings.Join(conditions, " OR ")
	query += " ORDER BY posts.id DESC;"

	// Execute the query
	rows, err := dbConn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var user User
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &user.Username, &user.LP)
		if err != nil {
			return nil, err
		}
		post.User = user
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPosts() ([]Post, error) {
	dbConn := db.GetDBConnection()

	rows, err := dbConn.Query(`
		SELECT posts.id, posts.title, posts.content, posts.category, posts.likes_count, posts.dislikes_count, posts.user_id, users.username, users.LP 
		FROM posts
		JOIN users ON posts.user_id = users.id
		ORDER BY posts.id DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var user User
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &user.Username, &user.LP)
		if err != nil {
			log.Printf("Error scanning post data: %v", err)
			return nil, err
		}
		post.User = user
		posts = append(posts, post)
	}

	return posts, nil
}

func PostDetailsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

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

	dbConn := db.GetDBConnection()
	var post Post
	err = dbConn.QueryRow(`
		SELECT posts.id, posts.title, posts.content, posts.category, posts.likes_count, posts.dislikes_count, posts.user_id, users.username, users.LP 
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.id = ?`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &post.User.Username, &post.User.LP)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	comments, err := FetchComments(post.ID)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Authenticated bool
		User          User
		Post          Post
		Comments      []Comment
	}{
		Authenticated: authenticated,
		User:          user,
		Post:          post,
		Comments:      comments,
	}

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
