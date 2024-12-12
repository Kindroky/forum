package handlers

import (
	"forum/db"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type HomepageData struct {
	Authenticated bool
	User          User
	Posts         []Post
}

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	LP        int
	SessionID string
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No session cookie found.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	dbConn := db.GetDBConnection()
	var user User
	err = dbConn.QueryRow(`
        SELECT id, username, LP, session_id 
        FROM users WHERE session_id = ?`, cookie.Value).Scan(&user.ID, &user.Username, &user.LP, &user.SessionID)
	if err != nil {
		log.Printf("Error retrieving user from session: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/addpost.html")
		if err != nil {
			log.Printf("Error parsing addpost template: %v", err)
			http.Error(w, "Error loading page", http.StatusInternalServerError)
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
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if len(categories) > 2 {
			http.Error(w, "You can only select up to 2 categories", http.StatusBadRequest)
			return
		}

		categoriesStr := strings.Join(categories, ",")

		err := db.CreatePost(title, content, categoriesStr, user.ID)
		if err != nil {
			log.Printf("Error creating post: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func toAnySlice(strings []string) []any {
	anySlice := make([]any, len(strings))
	for i, v := range strings {
		anySlice[i] = v
	}
	return anySlice
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
