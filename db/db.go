package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var database *sql.DB

// InitDB initializes the database connection and returns the connection object.
func InitDB(dataSourceName string) *sql.DB {
	var err error
	database, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	}
	CreateTables()
	// Ensure the session_id column exists
	ensureSessionIDColumn()

	return database
}

// GetDBConnection returns the current database connection.
func GetDBConnection() *sql.DB {
	return database
}

// Ensure the session_id column exists in the users table
func ensureSessionIDColumn() {
	// Check if session_id column exists, and add it if not
	_, err := database.Exec("ALTER TABLE users ADD COLUMN session_id TEXT;")
	if err != nil {
		if err.Error() != "duplicate column name: session_id" { // Ignore if the column already exists
			log.Printf("Error adding session_id column: %v", err)
		}
	}
}

// CreateTables initializes the database tables.
func CreateTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		session_id TEXT,
		LP INTEGER DEFAULT 0
	);`

	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		category TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	createLikesTable := `
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	createDislikesTable := `
	CREATE TABLE IF NOT EXISTS dislikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	createCommentTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err := database.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	_, err = database.Exec(createPostsTable)
	if err != nil {
		log.Fatalf("Error creating posts table: %v", err)
	}

	_, err = database.Exec(createLikesTable)
	if err != nil {
		log.Fatalf("Error creating likes table: %v", err)
	}

	_, err = database.Exec(createDislikesTable)
	if err != nil {
		log.Fatalf("Error creating dislikes table: %v", err)
	}

	_, err = database.Exec(createCommentTable)
	if err != nil {
		log.Fatalf("Error creating comment table: %v", err)
	}

	log.Println("Database tables initialized successfully.")
}

// CreatePost inserts a new post into the posts table.
func CreatePost(title, content, category string, userID int) error {
	_, err := database.Exec("INSERT INTO posts (title, content, category, user_id) VALUES (?, ?, ?, ?)", title, content, category, userID)
	if err != nil {
		log.Printf("Erreur lors de l'insertion du post : %v", err)
		return err
	}
	return nil
}

// GetPosts retrieves all posts from the database.
func GetPosts() ([]Post, error) {
	rows, err := database.Query(`
		SELECT 
			posts.id, posts.title, posts.content, posts.category, posts.user_id,
			users.id, users.username, 
			CASE 
				WHEN users.LP >= 100 THEN 'Legend'
				WHEN users.LP >= 50 THEN 'Pro'
				ELSE 'Novice'
			END AS rank, users.LP, users.session_id
		FROM posts
		JOIN users ON posts.user_id = users.id
		ORDER BY posts.id DESC;
	`)
	if err != nil {
		log.Printf("Erreur lors de la récupération des posts : %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var user User
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID,
			&user.ID, &user.Username, &user.Rank, &user.LP, &user.SessionID); err != nil {
			log.Printf("Erreur lors du scan des posts : %v", err)
			return nil, err
		}
		post.User = user
		posts = append(posts, post)
	}
	return posts, nil
}

// GetPostsByCategory retrieves posts filtered by category.
func GetPostsByCategory(category string) ([]Post, error) {
	rows, err := database.Query(`
		SELECT 
			posts.id, posts.title, posts.content, posts.category, posts.user_id,
			users.id, users.username, 
			CASE 
				WHEN users.LP >= 100 THEN 'Legend'
				WHEN users.LP >= 50 THEN 'Pro'
				ELSE 'Novice'
			END AS rank, users.LP, users.session_id
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.category = ?
		ORDER BY posts.id DESC;
	`, category)
	if err != nil {
		log.Printf("Erreur lors de la récupération des posts par catégorie : %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var user User
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID,
			&user.ID, &user.Username, &user.Rank, &user.LP, &user.SessionID); err != nil {
			log.Printf("Erreur lors du scan des posts : %v", err)
			return nil, err
		}
		post.User = user
		posts = append(posts, post)
	}
	return posts, nil
}

func GetComments(db *sql.DB, postID int) ([]Comment, error) {
	query := `SELECT comments.id, comments.content, comments.created_at, users.username 
              FROM comments 
              INNER JOIN users ON comments.user_id = users.id 
              WHERE comments.post_id = ? 
              ORDER BY comments.created_at ASC`

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// Struct for comments
type Comment struct {
	ID        int
	Content   string
	CreatedAt string
	Username  string
}

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
