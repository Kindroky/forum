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
	return database
}

// GetDBConnection returns the current database connection.
func GetDBConnection() *sql.DB {
	return database
}

func CreateTables() {
	// Create posts table if it doesn't exist
	createPostsTable := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        category TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id)
    );`
	_, err := database.Exec(createPostsTable)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table posts : %v", err)
	}

	// Ensure the "category" column exists in the posts table
	_, err = database.Exec("ALTER TABLE posts ADD COLUMN category TEXT NOT NULL DEFAULT 'general'")
	if err != nil && err.Error() != "duplicate column name: category" {
		log.Printf("Erreur lors de l'ajout de la colonne category : %v", err)
	}

	// Create users table if it doesn't exist
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        username TEXT NOT NULL,
        password TEXT NOT NULL
    );`
	_, err = database.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table users : %v", err)
	}

	log.Println("Tables créées ou déjà existantes.")
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
	rows, err := database.Query("SELECT id, title, content, category, user_id FROM posts")
	if err != nil {
		log.Printf("Erreur lors de la récupération des posts : %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID); err != nil {
			log.Printf("Erreur lors du scan des posts : %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetPostsByCategory retrieves posts filtered by category.
func GetPostsByCategory(category string) ([]Post, error) {
	rows, err := database.Query("SELECT id, title, content, category, user_id FROM posts WHERE category = ?", category)
	if err != nil {
		log.Printf("Erreur lors de la récupération des posts par catégorie : %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID); err != nil {
			log.Printf("Erreur lors du scan des posts : %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

type Post struct {
	ID       int
	Title    string
	Content  string
	Category string
	UserID   int
}
