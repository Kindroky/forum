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

type Post struct {
	ID       int
	Title    string
	Content  string
	Category string
	UserID   int
	User     User
}

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	Rank      string
	LP        int
	SessionID string
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No session cookie found.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	dbConn := db.GetDBConnection()
	var user db.User
	err = dbConn.QueryRow(`
		SELECT id, username 
		FROM users WHERE session_id = ?`, cookie.Value).Scan(&user.ID, &user.Username)
	if err != nil {
		log.Printf("Error retrieving user from session: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Handle GET request to display the form
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

	// Handle POST request to create a new post
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		// Validate form inputs
		if title == "" || content == "" || category == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Insert post into the database
		err := db.CreatePost(title, content, category, user.ID)
		if err != nil {
			log.Printf("Error creating post: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Redirect to the homepage after successful post creation
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user db.User

	// Check session
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
	
	// Fetch posts
	posts, err := db.GetPosts()
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create data for the template
	data := db.HomepageData{
		Authenticated: authenticated,
		User:          user,
		Posts:         posts,
	}

	// Render the template
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

	// Query to join posts and users table
	rows, err := dbConn.Query(`
		SELECT 
			posts.id, 
			posts.title, 
			posts.content, 
			posts.category, 
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
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID, &user.Username, &user.Rank, &user.LP)
		if err != nil {
			log.Printf("Error scanning post data: %v", err)
			return nil, err
		}
		post.User = user
		posts = append(posts, post)
	}

	return posts, nil
}
